# SVG Converter Pro 🎨

A powerful GitHub Action that converts SVG files to multiple formats including ICO, PNG, React components, and React Native components with extensive customization options.

## ✨ Features

- **Multi-format support**: ICO, PNG, React JS, React Native JS
- **Configurable output**: Custom sizes, TypeScript support, multiple PNGs
- **Smart naming**: Auto-detects names or uses custom base names
- **Professional logging**: Colored output with clear progress indicators
- **Comprehensive outputs**: JSON array of created files and human-readable summary

## 🚀 Quick Start

```yaml
- name: Convert SVG to multiple formats
  uses: ./.github/actions/svg-converter
  with:
    svg-path: 'assets/logo.svg'
    output-dir: 'dist/'
    formats: 'ico,png,react,react-native'
```

## 📋 Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `svg-path` | Path to the SVG file to convert | ✅ Yes | - |
| `output-dir` | Directory to output converted files | ❌ No | `.` |
| `formats` | Comma-separated formats: `ico,png,react,react-native` | ❌ No | `ico,png,react,react-native` |
| `png-sizes` | Comma-separated PNG sizes (e.g., `16,32,64,128,256`) | ❌ No | `16,32,64,128,256` |
| `ico-sizes` | Comma-separated ICO sizes (e.g., `16,32,48,64`) | ❌ No | `16,32,48,64` |
| `base-name` | Base name for output files (without extension) | ❌ No | *SVG filename* |
| `react-typescript` | Generate TypeScript React components | ❌ No | `false` |
| `react-props-interface` | Interface name for React component props | ❌ No | `SVGProps` |

## 📤 Outputs

| Output | Description |
|--------|-------------|
| `files-created` | JSON array of all created file paths |
| `summary` | Human-readable summary of conversion results |

## 🎯 Usage Examples

### Basic Conversion

```yaml
- name: Convert logo to all formats
  uses: ./.github/actions/svg-converter
  with:
    svg-path: 'branding/logo.svg'
    output-dir: 'assets/'
```

### Custom PNG Sizes

```yaml
- name: Convert with custom PNG sizes
  uses: ./.github/actions/svg-converter
  with:
    svg-path: 'icons/icon.svg'
    formats: 'png,ico'
    png-sizes: '16,24,32,48,64,96,128'
    ico-sizes: '16,32,48'
```

### TypeScript React Components

```yaml
- name: Generate TypeScript React components
  uses: ./.github/actions/svg-converter
  with:
    svg-path: 'icons/user.svg'
    formats: 'react,react-native'
    react-typescript: 'true'
    react-props-interface: 'IconProps'
    base-name: 'UserIcon'
```

### Favicon Generation

```yaml
- name: Generate favicon
  uses: ./.github/actions/svg-converter
  with:
    svg-path: 'branding/favicon.svg'
    formats: 'ico'
    ico-sizes: '16,32,48,64,128,256'
    base-name: 'favicon'
```

### Complete Workflow Example

```yaml
name: 🎨 Convert Assets

on:
  push:
    paths:
      - 'assets/**/*.svg'
  workflow_dispatch:

jobs:
  convert-assets:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Convert SVG Assets
        id: convert
        uses: ./.github/actions/svg-converter
        with:
          svg-path: 'assets/logo.svg'
          output-dir: 'public/assets/'
          formats: 'ico,png,react,react-native'
          png-sizes: '16,32,64,128,256,512'
          react-typescript: 'true'
          base-name: 'AppLogo'
      
      - name: Show conversion results
        run: |
          echo "Files created: ${{ steps.convert.outputs.files-created }}"
          echo "Summary:"
          echo "${{ steps.convert.outputs.summary }}"
      
      - name: Commit generated assets
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add public/assets/
          git diff --staged --quiet || git commit -m "🎨 Auto-generate assets from SVG"
          git push
```

## 🎨 Output Examples

### ICO Files

- `logo.ico` - Multi-resolution favicon with all specified sizes

### PNG Files

- `logo_16x16.png`
- `logo_32x32.png`
- `logo_64x64.png`
- `logo_128x128.png`
- `logo_256x256.png`

### React Components

- `Logo.js` or `Logo.tsx` - React component
- `Logo.native.js` or `Logo.native.tsx` - React Native component

## 🔧 Technical Details

### Dependencies

- **librsvg**: SVG to raster conversion
- **ImageMagick**: ICO format generation
- **@svgr/cli**: React component generation
- **jq**: JSON processing

### Supported Formats

- **ICO**: Multi-resolution favicon files
- **PNG**: Raster images at specified sizes
- **React**: JSX/TSX components for web
- **React Native**: JSX/TSX components for mobile

### Error Handling

- Validates SVG file existence and format
- Checks for required dependencies
- Provides detailed error messages
- Graceful handling of missing inputs

## 🎉 Why This Action Rocks

1. **One Action, Many Formats**: Convert once, get everything you need
2. **Smart Defaults**: Works great out of the box
3. **Highly Configurable**: Customize every aspect
4. **TypeScript Ready**: Full TypeScript support for React components  
5. **Professional Output**: Beautiful logging and comprehensive results
6. **Production Ready**: Robust error handling and validation

Perfect for maintaining consistent branding assets across web, mobile, and favicon use cases! 🚀
