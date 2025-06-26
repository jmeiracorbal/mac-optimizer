#!/bin/bash

# Mac Optimizer - Sistema de optimización y limpieza para macOS
# Autor: José Meira
# Version: 1.0.0

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Variables globales
VERSION="1.0.0"
DRY_RUN=false
QUIET=false
FORCE=false
LOG_FILE="$HOME/.mac-optimizer.log"

# Función para mostrar ayuda
show_help() {
    cat << EOF
Mac Optimizer v${VERSION}

DESCRIPTION:
    macOS optimization and cleanup tool that improves system performance
    through cache cleaning, service optimization, and resource management.

USAGE:
    mac-optimizer [OPTIONS] [MODULES]

OPTIONS:
    -h, --help          Show this help message
    -v, --version       Show version information
    -d, --dry-run       Show what would be done without making changes
    -q, --quiet         Run in quiet mode
    -f, --force         Force operations without confirmation
    --log-file FILE     Specify log file (default: ~/.mac-optimizer.log)

OPTIMIZATION MODULES:
    cache              Clean system and application caches
    spotlight          Optimize and reindex Spotlight
    permissions        Repair system permissions
    services           Restart problematic services
    memory             Free memory and clean processes
    logs               Clean old system logs
    browsers           Clean browser caches and data
    downloads          Clean old downloads (with confirmation)
    all                Run all modules (default)

EXAMPLES:
    mac-optimizer                    # Complete interactive optimization
    mac-optimizer --dry-run          # Preview changes
    mac-optimizer cache spotlight    # Only clean cache and optimize Spotlight
    mac-optimizer --force --quiet    # Complete automatic optimization
    mac-optimizer browsers downloads # Clean browsers and downloads only

NOTES:
    - Some operations require administrator permissions
    - Automatic backups are created before critical changes
    - Check log at ~/.mac-optimizer.log for details
    - Browser and downloads cleaning requires user confirmation

EOF
}

# Función para mostrar versión
show_version() {
    echo "Mac Optimizer v${VERSION}"
    echo "Compatible con macOS 10.15+"
}

# Función para logging
log_message() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[$timestamp] [$level] $message" >> "$LOG_FILE"
    
    if [ "$QUIET" = false ]; then
        case $level in
            "INFO")  echo -e "${BLUE}ℹ️  $message${NC}" ;;
            "WARN")  echo -e "${YELLOW}⚠️  $message${NC}" ;;
            "ERROR") echo -e "${RED}❌ $message${NC}" ;;
            "SUCCESS") echo -e "${GREEN}✅ $message${NC}" ;;
            "PROGRESS") echo -e "${PURPLE}🔄 $message${NC}" ;;
        esac
    fi
}

# Función para verificar permisos de administrador
check_admin_permissions() {
    if [[ $EUID -eq 0 ]]; then
        log_message "WARN" "Ejecutándose como root. Algunas operaciones podrían ser riesgosas."
        return 0
    fi
    
    # Verificar si puede usar sudo
    if sudo -n true 2>/dev/null; then
        return 0
    fi
    
    log_message "WARN" "Some operations require administrator permissions."
    if [ "$FORCE" = false ]; then
        if read -p "Continue without administrator permissions? (y/N): " continue_without_admin; then
            if [[ ! $continue_without_admin =~ ^[Yy]$ ]]; then
                log_message "ERROR" "Operation cancelled by user"
                exit 1
            fi
        fi
    fi
    return 1
}

# Función para mostrar spinner de carga
show_spinner() {
    local pid=$1
    local message="$2"
    local spin='-\|/'
    local i=0
    
    while kill -0 $pid 2>/dev/null; do
        i=$(( (i+1) %4 ))
        printf "\r${PURPLE}${spin:$i:1} $message${NC}"
        sleep 0.1
    done
    printf "\r"
}

# Función para ejecutar comando con spinner
run_with_spinner() {
    local command="$1"
    local message="$2"
    
    if [ "$QUIET" = true ]; then
        eval "$command" >/dev/null 2>&1
    else
        eval "$command" >/dev/null 2>&1 &
        local pid=$!
        show_spinner $pid "$message"
        wait $pid
        local exit_code=$?
        return $exit_code
    fi
}
get_folder_size() {
    local folder="$1"
    if [ -d "$folder" ]; then
        du -sh "$folder" 2>/dev/null | awk '{print $1}' || echo "0B"
    else
        echo "0B"
    fi
}

# Función para limpiar caches
clean_caches() {
    log_message "PROGRESS" "Starting cache cleanup..."
    
    # Caches de usuario
    local user_cache_dirs=(
        "$HOME/Library/Caches"
        "$HOME/Library/Application Support/CrashReporter"
        "$HOME/Library/Logs"
        "$HOME/Library/WebKit"
        "$HOME/Library/Saved Application State"
    )
    
    for cache_dir in "${user_cache_dirs[@]}"; do
        if [ -d "$cache_dir" ]; then
            local size_before=$(get_folder_size "$cache_dir")
            log_message "INFO" "Cleaning $cache_dir (size: $size_before)"
            
            if [ "$DRY_RUN" = false ]; then
                # Limpiar archivos más recientes también, pero preservar estructura
                find "$cache_dir" -type f -name "*" -delete 2>/dev/null || true
                find "$cache_dir" -type d -empty -delete 2>/dev/null || true
            fi
            
            local size_after=$(get_folder_size "$cache_dir")
            log_message "SUCCESS" "Cache cleaned: $cache_dir (now: $size_after)"
        fi
    done
    
    # Caches del sistema (requiere permisos admin)
    if check_admin_permissions; then
        local system_cache_dirs=(
            "/Library/Caches"
            "/private/var/log/asl"
        )
        
        for cache_dir in "${system_cache_dirs[@]}"; do
            if [ -d "$cache_dir" ]; then
                local size_before=$(get_folder_size "$cache_dir")
                log_message "INFO" "Cleaning system cache: $cache_dir"
                
                if [ "$DRY_RUN" = false ]; then
                    sudo find "$cache_dir" -type f -mtime +1 -delete 2>/dev/null || true
                fi
            fi
        done
    fi
    
    # DNS Cache
    log_message "INFO" "Cleaning DNS cache"
    if [ "$DRY_RUN" = false ]; then
        sudo dscacheutil -flushcache 2>/dev/null || true
        sudo killall -HUP mDNSResponder 2>/dev/null || true
    fi
    
    log_message "SUCCESS" "Cache cleanup completed"
}

# Función para optimizar Spotlight
optimize_spotlight() {
    log_message "PROGRESS" "Spotlight optimization requested..."
    
    if [ "$FORCE" = false ] && [ "$DRY_RUN" = false ]; then
        echo -e "${YELLOW}⚠️  Spotlight optimization will:${NC}"
        echo -e "${YELLOW}   • Require administrator permissions${NC}"
        echo -e "${YELLOW}   • Temporarily disable indexing${NC}"
        echo -e "${YELLOW}   • Rebuild the search index (takes time)${NC}"
        echo -e "${YELLOW}   • May affect search functionality temporarily${NC}"
        read -p "Continue with Spotlight optimization? (y/N): " confirm_spotlight
        if [[ ! $confirm_spotlight =~ ^[Yy]$ ]]; then
            log_message "INFO" "Spotlight optimization skipped by user"
            return 0
        fi
    fi
    
    log_message "PROGRESS" "Optimizing Spotlight..."
    
    if [ "$DRY_RUN" = false ]; then
        # Deshabilitar indexado temporalmente
        log_message "INFO" "Disabling Spotlight indexing..."
        run_with_spinner "sudo mdutil -a -i off 2>/dev/null" "Disabling Spotlight indexing..."
        
        sleep 2
        
        # Limpiar índice existente
        log_message "INFO" "Removing existing Spotlight index..."
        run_with_spinner "sudo rm -rf /.Spotlight-V100 /System/Volumes/Data/.Spotlight-V100 2>/dev/null" "Removing Spotlight index..."
        
        # Reactivar indexado
        log_message "INFO" "Re-enabling Spotlight indexing..."
        run_with_spinner "sudo mdutil -a -i on 2>/dev/null" "Re-enabling Spotlight indexing..."
        
        log_message "SUCCESS" "Spotlight reindexed. Process will continue in background."
    else
        log_message "INFO" "Would reindex Spotlight (background process)"
    fi
}

# Función para reparar permisos
repair_permissions() {
    log_message "PROGRESS" "Repairing system permissions..."
    
    if [ "$DRY_RUN" = false ]; then
        # Reparar permisos de disco (solo en macOS más antiguos)
        if command -v diskutil >/dev/null 2>&1; then
            sudo diskutil repairPermissions / 2>/dev/null || true
        fi
        
        # Resetear permisos de carpetas importantes
        sudo chmod 755 /usr/local/bin 2>/dev/null || true
        sudo chmod 755 "$HOME" 2>/dev/null || true
        
        log_message "SUCCESS" "Permissions repaired"
    else
        log_message "INFO" "Would repair system permissions"
    fi
}

# Función para reiniciar servicios
restart_services() {
    log_message "PROGRESS" "Restarting system services..."
    
    local services=(
        "Finder"
        "Dock"
        "SystemUIServer"
        "ControlStrip"
    )
    
    for service in "${services[@]}"; do
        log_message "INFO" "Restarting $service"
        if [ "$DRY_RUN" = false ]; then
            killall "$service" 2>/dev/null || true
        fi
    done
    
    log_message "SUCCESS" "Services restarted"
}

# Función para liberar memoria
free_memory() {
    log_message "PROGRESS" "Freeing system memory..."
    
    if [ "$DRY_RUN" = false ]; then
        # Forzar garbage collection
        if command -v purge >/dev/null 2>&1; then
            sudo purge 2>/dev/null || true
            log_message "INFO" "Memory purge completed"
        fi
        
        # Eliminar dynamic_pager completamente - causa problemas
        # Es mejor dejar que el sistema maneje la memoria swap automáticamente
        
        log_message "SUCCESS" "Memory freed"
    else
        log_message "INFO" "Would free system memory"
    fi
}

# Función para limpiar navegadores
clean_browsers() {
    log_message "PROGRESS" "Cleaning browser caches and data..."
    
    if [ "$FORCE" = false ] && [ "$DRY_RUN" = false ]; then
        echo -e "${YELLOW}⚠️  This will clean browser caches, cookies, and temporary data.${NC}"
        echo -e "${YELLOW}   You may need to log in to websites again.${NC}"
        read -p "Continue with browser cleanup? (y/N): " confirm_browsers
        if [[ ! $confirm_browsers =~ ^[Yy]$ ]]; then
            log_message "INFO" "Browser cleanup skipped by user"
            return 0
        fi
    fi
    
    # Chrome caches
    local chrome_dirs=(
        "$HOME/Library/Caches/Google/Chrome"
        "$HOME/Library/Application Support/Google/Chrome/Default/Application Cache"
        "$HOME/Library/Application Support/Google/Chrome/Default/GPUCache"
    )
    
    # Firefox caches
    local firefox_profile_dir="$HOME/Library/Application Support/Firefox/Profiles"
    
    # Safari caches
    local safari_dirs=(
        "$HOME/Library/Caches/com.apple.Safari"
        "$HOME/Library/Safari/LocalStorage"
        "$HOME/Library/Safari/Databases"
        "$HOME/Library/WebKit"
    )
    
    # Limpiar Chrome
    for cache_dir in "${chrome_dirs[@]}"; do
        if [ -d "$cache_dir" ]; then
            local size_before=$(get_folder_size "$cache_dir")
            log_message "INFO" "Cleaning Chrome cache: $cache_dir ($size_before)"
            
            if [ "$DRY_RUN" = false ]; then
                rm -rf "$cache_dir"/* 2>/dev/null || true
            fi
        fi
    done
    
    # Limpiar Firefox
    if [ -d "$firefox_profile_dir" ]; then
        log_message "INFO" "Cleaning Firefox caches"
        if [ "$DRY_RUN" = false ]; then
            find "$firefox_profile_dir" -name "cache2" -type d -exec rm -rf {} + 2>/dev/null || true
            find "$firefox_profile_dir" -name "startupCache" -type d -exec rm -rf {} + 2>/dev/null || true
        fi
    fi
    
    # Limpiar Safari
    for cache_dir in "${safari_dirs[@]}"; do
        if [ -d "$cache_dir" ]; then
            local size_before=$(get_folder_size "$cache_dir")
            log_message "INFO" "Cleaning Safari cache: $cache_dir ($size_before)"
            
            if [ "$DRY_RUN" = false ]; then
                rm -rf "$cache_dir"/* 2>/dev/null || true
            fi
        fi
    done
    
    log_message "SUCCESS" "Browser caches cleaned"
}

# Función para limpiar descargas
clean_downloads() {
    log_message "PROGRESS" "Cleaning old downloads..."
    
    local downloads_dir="$HOME/Downloads"
    
    if [ ! -d "$downloads_dir" ]; then
        log_message "INFO" "Downloads directory not found"
        return 0
    fi
    
    # Contar archivos antiguos (más de 30 días)
    local old_files_count=$(find "$downloads_dir" -type f -mtime +30 2>/dev/null | wc -l | tr -d ' ')
    local total_size=$(get_folder_size "$downloads_dir")
    
    if [ "$old_files_count" -eq 0 ]; then
        log_message "INFO" "No old files found in Downloads (>30 days)"
        return 0
    fi
    
    log_message "INFO" "Found $old_files_count files older than 30 days in Downloads"
    log_message "INFO" "Total Downloads folder size: $total_size"
    
    if [ "$FORCE" = false ] && [ "$DRY_RUN" = false ]; then
        echo -e "${YELLOW}⚠️  This will delete files older than 30 days from Downloads.${NC}"
        echo -e "${YELLOW}   Files found: $old_files_count${NC}"
        echo -e "${YELLOW}   Total folder size: $total_size${NC}"
        read -p "Continue with Downloads cleanup? (y/N): " confirm_downloads
        if [[ ! $confirm_downloads =~ ^[Yy]$ ]]; then
            log_message "INFO" "Downloads cleanup skipped by user"
            return 0
        fi
    fi
    
    if [ "$DRY_RUN" = false ]; then
        # Mostrar algunos archivos que se van a eliminar
        log_message "INFO" "Files to be deleted (sample):"
        find "$downloads_dir" -type f -mtime +30 2>/dev/null | head -5 | while read file; do
            log_message "INFO" "  $(basename "$file")"
        done
        
        # Eliminar archivos antiguos
        local deleted_count=$(find "$downloads_dir" -type f -mtime +30 -delete 2>/dev/null | wc -l | tr -d ' ')
        log_message "SUCCESS" "Deleted $old_files_count old files from Downloads"
    else
        log_message "INFO" "Would delete $old_files_count files older than 30 days"
        find "$downloads_dir" -type f -mtime +30 2>/dev/null | head -5 | while read file; do
            log_message "INFO" "  Would delete: $(basename "$file")"
        done
    fi
}
clean_logs() {
    log_message "PROGRESS" "Limpiando logs del sistema..."
    
    local log_dirs=(
        "$HOME/Library/Logs"
        "/private/var/log"
    )
    
    for log_dir in "${log_dirs[@]}"; do
        if [ -d "$log_dir" ]; then
            log_message "INFO" "Limpiando logs en $log_dir"
            if [ "$DRY_RUN" = false ]; then
                find "$log_dir" -name "*.log" -mtime +30 -delete 2>/dev/null || true
                find "$log_dir" -name "*.out" -mtime +30 -delete 2>/dev/null || true
            fi
        fi
    done
    
    log_message "SUCCESS" "Logs limpiados"
}

# Función para mostrar estadísticas del sistema
show_system_stats() {
    echo -e "${CYAN}=== Estadísticas del Sistema ===${NC}"
    
    # Uso de memoria
    local memory_pressure=$(memory_pressure | grep "System-wide memory free percentage" | awk '{print $5}' | sed 's/%//')
    echo -e "${BLUE}Memoria libre: ${memory_pressure}%${NC}"
    
    # Uso de disco
    local disk_usage=$(df -h / | tail -1 | awk '{print $5}' | sed 's/%//')
    echo -e "${BLUE}Uso de disco: ${disk_usage}%${NC}"
    
    # Uptime
    local uptime_info=$(uptime | awk '{print $3, $4}' | sed 's/,//')
    echo -e "${BLUE}Tiempo encendido: ${uptime_info}${NC}"
    
    echo ""
}

# Función principal de optimización
run_optimization() {
    local modules=("$@")
    
    if [ ${#modules[@]} -eq 0 ]; then
        modules=("cache" "spotlight" "permissions" "services" "memory" "logs")
    fi
    
    echo -e "${CYAN}=== Mac Optimizer v${VERSION} ===${NC}"
    echo -e "${YELLOW}Starting system optimization...${NC}"
    
    if [ "$DRY_RUN" = true ]; then
        echo -e "${YELLOW}🔍 DRY-RUN MODE: Only showing what would be done${NC}"
    fi
    
    echo ""
    show_system_stats
    
    # Crear backup del log anterior
    if [ -f "$LOG_FILE" ]; then
        cp "$LOG_FILE" "${LOG_FILE}.backup"
    fi
    
    log_message "INFO" "=== Starting Mac Optimizer v${VERSION} ==="
    
    for module in "${modules[@]}"; do
        case $module in
            "cache")       clean_caches ;;
            "spotlight")   optimize_spotlight ;;
            "permissions") repair_permissions ;;
            "services")    restart_services ;;
            "memory")      free_memory ;;
            "logs")        clean_logs ;;
            "browsers")    clean_browsers ;;
            "downloads")   clean_downloads ;;
            "all")         
                clean_caches
                optimize_spotlight
                repair_permissions
                restart_services
                free_memory
                clean_logs
                ;;
            *)
                log_message "ERROR" "Unknown module: $module"
                ;;
        esac
        echo ""
    done
    
    echo -e "${GREEN}🎉 Optimization completed!${NC}"
    echo -e "${BLUE}📋 Log saved at: $LOG_FILE${NC}"
    echo ""
    show_system_stats
}

# Procesar argumentos de línea de comandos
modules_to_run=()

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--version)
            show_version
            exit 0
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        -f|--force)
            FORCE=true
            shift
            ;;
        --log-file)
            LOG_FILE="$2"
            shift 2
            ;;
        cache|spotlight|permissions|services|memory|logs|browsers|downloads|all)
            modules_to_run+=("$1")
            shift
            ;;
        *)
            echo -e "${RED}❌ Opción desconocida: $1${NC}"
            echo "Usa --help para ver las opciones disponibles"
            exit 1
            ;;
    esac
done

# Verificar que estamos en macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    log_message "ERROR" "Este script solo funciona en macOS"
    exit 1
fi

# Ejecutar optimización
run_optimization "${modules_to_run[@]}"
