# Poo Tracker Branding Assets

This directory contains all the branding assets for the Poo Tracker application, generated from the source `logo.svg` file using our automated SVG conversion workflow.

## Generated Assets

### Icons & Favicons

- `favicon.ico` - Browser favicon (16x16, 32x32, 48x48 multi-resolution)
- `logo_16x16.png` - Small icon for UI elements
- `logo_32x32.png` - Standard icon size
- `logo_64x64.png` - Medium icon for cards/buttons
- `logo_128x128.png` - Large icon for headers
- `logo_256x256.png` - High-res icon for PWA/mobile

### Vector Graphics

- `logo.svg` - Source vector logo (scalable)

### React Components

- `logo.js` - React component (JavaScript)
- `logo.native.js` - React Native component

### Source Files

- `logo.svg` - Original design file
- `poop-emoji.ai` - Adobe Illustrator source

## Usage in Application

### Frontend Integration

The branding assets are integrated throughout the poo-tracker frontend:

1. **HTML Head** (`frontend/index.html`):

   - Multiple favicon sizes for different contexts
   - PWA manifest integration
   - Mobile app meta tags

2. **React Components**:

   - `Logo.tsx` - TypeScript component with size/className props
   - Used in Navbar, HomePage, and LoginPage

3. **PWA Manifest** (`frontend/public/manifest.json`):
   - App name, description, theme colors
   - Icon definitions for installation
   - Standalone app configuration

### Backend Integration

The backend serves branding assets via static file middleware:

- Endpoint: `/assets/logo.svg`, `/assets/logo_*.png`
- Used for API responses and email templates

## Conversion Workflow

Assets are generated using our custom SVG Converter GitHub Action:

- **Trigger**: Push to `main` or manual workflow dispatch
- **Input**: `branding/logo.svg`
- **Output**: All asset formats automatically generated
- **Quality**: Inkscape-powered conversion for optimal rendering

## Asset Guidelines

### Usage Standards

1. **Primary Logo**: Use `logo.svg` for all scalable contexts
2. **Favicons**: Automatically handled by HTML head configuration
3. **React Components**: Import `Logo` component, specify size prop
4. **Mobile Apps**: Use appropriate PNG sizes for different screen densities

### Brand Colors

- **Primary Brown**: `#8B4513` (theme-color)
- **Secondary Brown**: `#B35F37` (logo accent)
- **Text Brown**: `#6B4423` (readable text)

### Size Recommendations

- **Navbar**: 32px
- **Hero Section**: 96px
- **Login Form**: 64px
- **Mobile Favicon**: 32x32px
- **Desktop Favicon**: 16x16px
- **App Icon**: 256x256px

## Maintenance

### Updating the Logo

1. Edit `branding/logo.svg` with your design changes
2. Commit and push to trigger automatic asset generation
3. New assets will be generated and committed automatically
4. Update any hardcoded references if needed

### Adding New Formats

Modify `.github/workflows/svg-convert.yml` to include additional output formats or sizes as needed.

## Quality Assurance

All generated assets maintain:

- ✅ Consistent visual identity
- ✅ Proper transparency handling
- ✅ Optimal file sizes
- ✅ Cross-browser compatibility
- ✅ PWA compliance
- ✅ Mobile-first responsive design
