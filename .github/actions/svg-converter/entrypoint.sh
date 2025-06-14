#!/bin/bash

# SVG Converter Pro - Entrypoint Script
#
# This script converts SVG files to multiple formats with configurable options:
#   - ICO (multi-resolution favicons)
#   - PNG (various sizes)
#   - React JS components
#   - React Native JS components
#
# DEPENDENCIES:
#   - librsvg (for rsvg-convert)
#   - imagemagick (for convert)
#   - @svgr/cli (for React components)
#   - jq (for JSON processing)

set -euo pipefail

# Colors for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly BOLD='\033[1m'
readonly NC='\033[0m' # No Color

# Input variables from GitHub Actions
readonly SVG_PATH="${INPUT_SVG_PATH}"
readonly OUTPUT_DIR="${INPUT_OUTPUT_DIR:-./}"
readonly FORMATS="${INPUT_FORMATS:-ico,png,react,react-native}"
readonly PNG_SIZES="${INPUT_PNG_SIZES:-16,32,64,128,256}"
readonly ICO_SIZES="${INPUT_ICO_SIZES:-16,32,48,64}"
readonly BASE_NAME="${INPUT_BASE_NAME:-}"
readonly REACT_TYPESCRIPT="${INPUT_REACT_TYPESCRIPT:-false}"
readonly REACT_PROPS_INTERFACE="${INPUT_REACT_PROPS_INTERFACE:-SVGProps}"

# Global variables
declare -a CREATED_FILES=()
declare -a CONVERSION_SUMMARY=()

# Logging functions
log_info() {
    echo -e "${GREEN}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}" >&2
}

log_step() {
    echo -e "${BLUE}${BOLD}🔄 $1${NC}"
}

# Check if required tools are available
check_dependencies() {
    local missing_deps=()

    command -v rsvg-convert >/dev/null 2>&1 || missing_deps+=("librsvg")
    command -v convert >/dev/null 2>&1 || missing_deps+=("imagemagick")
    command -v svgr >/dev/null 2>&1 || missing_deps+=("@svgr/cli")
    command -v jq >/dev/null 2>&1 || missing_deps+=("jq")

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        exit 1
    fi
}

# Validate inputs
validate_inputs() {
    if [[ ! -f "$SVG_PATH" ]]; then
        log_error "SVG file not found: $SVG_PATH"
        exit 1
    fi

    if [[ ! "$SVG_PATH" =~ \.svg$ ]]; then
        log_error "File is not an SVG: $SVG_PATH"
        exit 1
    fi

    # Create output directory if it doesn't exist
    mkdir -p "$OUTPUT_DIR"
}

# Get base name for output files
get_base_name() {
    if [[ -n "$BASE_NAME" ]]; then
        echo "$BASE_NAME"
    else
        basename "$SVG_PATH" .svg
    fi
}

# Convert SVG to ICO format
convert_to_ico() {
    local base_name="$1"
    local output_file="$OUTPUT_DIR/${base_name}.ico"
    local tmp_png="/tmp/${base_name}_temp.png"

    log_step "Converting to ICO format..."

    # Convert SVG to high-res PNG first
    rsvg-convert -w 256 -h 256 "$SVG_PATH" -o "$tmp_png"

    # Create multi-resolution ICO
    convert "$tmp_png" -define icon:auto-resize="$ICO_SIZES" "$output_file"

    # Cleanup
    rm -f "$tmp_png"

    CREATED_FILES+=("$output_file")
    CONVERSION_SUMMARY+=("ICO: $output_file (sizes: $ICO_SIZES)")
    log_success "Created ICO: $output_file"
}

# Convert SVG to PNG format(s)
convert_to_png() {
    local base_name="$1"

    log_step "Converting to PNG format(s)..."

    IFS=',' read -ra SIZES <<<"$PNG_SIZES"
    for size in "${SIZES[@]}"; do
        local output_file="$OUTPUT_DIR/${base_name}_${size}x${size}.png"

        rsvg-convert -w "$size" -h "$size" "$SVG_PATH" -o "$output_file"

        CREATED_FILES+=("$output_file")
        log_success "Created PNG: $output_file"
    done

    CONVERSION_SUMMARY+=("PNG: ${#SIZES[@]} files created (sizes: $PNG_SIZES)")
}

# Convert SVG to React component
convert_to_react() {
    local base_name="$1"
    local extension="js"
    local output_file

    if [[ "$REACT_TYPESCRIPT" == "true" ]]; then
        extension="tsx"
    fi

    output_file="$OUTPUT_DIR/${base_name}.${extension}"

    log_step "Converting to React component..."

    local svgr_args=()

    # Add TypeScript flag if requested
    if [[ "$REACT_TYPESCRIPT" == "true" ]]; then
        svgr_args+=(--typescript)
    fi

    # Add props interface if specified
    if [[ -n "$REACT_PROPS_INTERFACE" ]]; then
        svgr_args+=(--props-interface "$REACT_PROPS_INTERFACE")
    fi

    # Convert SVG to React component
    svgr "${svgr_args[@]}" --out-file "$output_file" "$SVG_PATH"

    CREATED_FILES+=("$output_file")
    CONVERSION_SUMMARY+=("React: $output_file (TypeScript: $REACT_TYPESCRIPT)")
    log_success "Created React component: $output_file"
}

# Convert SVG to React Native component
convert_to_react_native() {
    local base_name="$1"
    local extension="js"
    local output_file

    if [[ "$REACT_TYPESCRIPT" == "true" ]]; then
        extension="tsx"
    fi

    output_file="$OUTPUT_DIR/${base_name}.native.${extension}"

    log_step "Converting to React Native component..."

    local svgr_args=(--native)

    # Add TypeScript flag if requested
    if [[ "$REACT_TYPESCRIPT" == "true" ]]; then
        svgr_args+=(--typescript)
    fi

    # Add props interface if specified
    if [[ -n "$REACT_PROPS_INTERFACE" ]]; then
        svgr_args+=(--props-interface "$REACT_PROPS_INTERFACE")
    fi

    # Convert SVG to React Native component
    svgr "${svgr_args[@]}" --out-file "$output_file" "$SVG_PATH"

    CREATED_FILES+=("$output_file")
    CONVERSION_SUMMARY+=("React Native: $output_file (TypeScript: $REACT_TYPESCRIPT)")
    log_success "Created React Native component: $output_file"
}

# Set GitHub Actions outputs
set_outputs() {
    local files_json
    local summary_text

    # Convert array to JSON
    files_json=$(printf '%s\n' "${CREATED_FILES[@]}" | jq -R . | jq -s .)

    # Create summary text
    summary_text=$(printf "Converted %s to %d files:\n%s" "$SVG_PATH" "${#CREATED_FILES[@]}" "$(printf '%s\n' "${CONVERSION_SUMMARY[@]}")")

    # Set outputs
    echo "files-created=$files_json" >>"$GITHUB_OUTPUT"
    echo "summary<<EOF" >>"$GITHUB_OUTPUT"
    echo "$summary_text" >>"$GITHUB_OUTPUT"
    echo "EOF" >>"$GITHUB_OUTPUT"
}

# Main conversion function
main() {
    log_info "🎨 SVG Converter Pro - Starting conversion..."
    log_info "📁 Input: $SVG_PATH"
    log_info "📁 Output: $OUTPUT_DIR"
    log_info "🎯 Formats: $FORMATS"

    check_dependencies
    validate_inputs

    local base_name
    base_name=$(get_base_name)
    log_info "📝 Base name: $base_name"

    # Parse requested formats
    IFS=',' read -ra FORMAT_ARRAY <<<"$FORMATS"

    for format in "${FORMAT_ARRAY[@]}"; do
        case "$format" in
        ico)
            convert_to_ico "$base_name"
            ;;
        png)
            convert_to_png "$base_name"
            ;;
        react)
            convert_to_react "$base_name"
            ;;
        react-native)
            convert_to_react_native "$base_name"
            ;;
        *)
            log_warn "Unknown format: $format"
            ;;
        esac
    done

    set_outputs

    log_success "🎉 Conversion completed! Created ${#CREATED_FILES[@]} files."

    # Print summary
    echo -e "\n${BOLD}📋 CONVERSION SUMMARY:${NC}"
    printf '%s\n' "${CONVERSION_SUMMARY[@]}"
}

# Run main function
main "$@"
