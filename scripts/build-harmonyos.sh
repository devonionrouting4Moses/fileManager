#!/bin/bash
set -e

# Get the project root (parent of scripts directory)
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# Read version from VERSION file
if [ -f "$PROJECT_ROOT/VERSION" ]; then
    VERSION=$(cat "$PROJECT_ROOT/VERSION" | tr -d ' \n')
else
    VERSION="2.0.0"
    echo "⚠️  VERSION file not found, using default: $VERSION"
fi

echo "🔨 Building FileManager v${VERSION} for HarmonyOS..."

# Check if OpenHarmony SDK is installed
if [ -z "$OHOS_SDK_HOME" ]; then
    echo "⚠️  OHOS_SDK_HOME not set. Please install OpenHarmony SDK first."
    echo "Download from: https://developer.harmonyos.com/"
    echo ""
    echo "After installation, set:"
    echo "  export OHOS_SDK_HOME=/path/to/ohos-sdk"
    exit 1
fi

echo "✅ OpenHarmony SDK found at: $OHOS_SDK_HOME"

# Create HAP structure
echo "📦 Creating HarmonyOS Ability Package (HAP) structure..."
mkdir -p harmonyos/entry/src/main/ets/pages
mkdir -p harmonyos/entry/src/main/ets/common
mkdir -p harmonyos/entry/src/main/resources/base/element
mkdir -p harmonyos/entry/src/main/resources/base/media
mkdir -p harmonyos/entry/src/main/resources/base/profile

# Create module.json5
cat > harmonyos/entry/src/main/module.json5 << 'EOFMODULE'
{
  "module": {
    "name": "entry",
    "type": "entry",
    "srcEntry": "./ets/Application/AbilityStage.ts",
    "description": "$string:module_desc",
    "mainElement": "MainAbility",
    "deviceTypes": ["phone", "tablet"],
    "abilities": [
      {
        "name": "MainAbility",
        "srcEntry": "./ets/pages/Index.ets",
        "description": "$string:MainAbility_desc",
        "icon": "$media:icon",
        "label": "$string:MainAbility_label",
        "startWindowIcon": "$media:icon",
        "startWindowBackground": "$color:start_window_background",
        "exported": true,
        "skills": [
          {
            "entities": ["entity.system.home"],
            "actions": ["action.system.home"]
          }
        ]
      }
    ]
  }
}
EOFMODULE

# Create AbilityStage.ts
cat > harmonyos/entry/src/main/ets/Application/AbilityStage.ts << 'EOFABILITY'
import AbilityStage from '@ohos.app.ability.AbilityStage';

export default class MyAbilityStage extends AbilityStage {
  onCreate() {
    console.log("FileManager AbilityStage onCreate");
  }
}
EOFABILITY

# Create Index.ets (main UI)
cat > harmonyos/entry/src/main/ets/pages/Index.ets << 'EOFINDEX'
import { router } from '@kit.ArkUI';

@Entry
@Component
struct Index {
  @State message: string = 'FileManager v2.0.0';

  build() {
    RelativeContainer() {
      Column() {
        Text(this.message)
          .fontSize(50)
          .fontWeight(FontWeight.Bold)
          .textAlign(TextAlign.Center)

        Button('Open File Manager')
          .width(200)
          .height(50)
          .margin({ top: 20 })
          .onClick(() => {
            // File manager functionality
          })
      }
      .width('100%')
      .height('100%')
      .justifyContent(FlexAlign.Center)
      .alignItems(HorizontalAlign.Center)
    }
    .height('100%')
  }
}
EOFINDEX

# Create package.json
cat > harmonyos/entry/package.json << 'EOFPACKAGE'
{
  "name": "filemanager",
  "version": "2.0.0",
  "description": "Modern file manager with Rust+Go backend",
  "main": "index.ets",
  "author": "DevChigarlicMoses",
  "license": "MIT",
  "dependencies": {
    "@ohos/hvigor": "^4.0.0"
  }
}
EOFPACKAGE

# Create build.gradle
cat > harmonyos/entry/build.gradle << 'EOFGRADLE'
apply plugin: 'com.huawei.ohos.app'

ohos {
    compileSdkVersion 12
    defaultConfig {
        compatibleSdkVersion 12
    }
    buildTypes {
        release {
            minifyEnabled true
            shrinkResources true
        }
    }
}

dependencies {
    implementation fileTree(dir: 'libs', include: ['*.jar', '*.aar'])
}
EOFGRADLE

# Create strings.json
cat > harmonyos/entry/src/main/resources/base/element/string.json << 'EOFSTRINGS'
{
  "string": [
    {
      "name": "module_desc",
      "value": "module description"
    },
    {
      "name": "MainAbility_desc",
      "value": "mainability description"
    },
    {
      "name": "MainAbility_label",
      "value": "FileManager"
    }
  ]
}
EOFSTRINGS

# Create colors.json
cat > harmonyos/entry/src/main/resources/base/element/color.json << 'EOFCOLORS'
{
  "color": [
    {
      "name": "start_window_background",
      "value": "#FFFFFF"
    }
  ]
}
EOFCOLORS

echo "✅ HAP structure created"

# Try to build HAP
echo "🏗️  Building HAP package..."
if command -v hvigorw &> /dev/null; then
    cd harmonyos
    hvigorw build --mode module
    cd ..
    echo "✅ HAP package built successfully"
    echo "Output: harmonyos/entry/build/outputs/hap/entry-default-signed.hap"
else
    echo "⚠️  hvigorw not found in PATH"
    echo "Please ensure OpenHarmony SDK is properly installed and configured"
    echo "HAP structure has been created in: harmonyos/"
    echo "You can build it manually with: cd harmonyos && hvigorw build --mode module"
fi

echo ""
echo "🎉 HarmonyOS package preparation completed!"
echo ""
echo "Files created:"
echo "  - harmonyos/entry/src/main/module.json5"
echo "  - harmonyos/entry/src/main/ets/pages/Index.ets"
echo "  - harmonyos/entry/src/main/ets/Application/AbilityStage.ts"
echo "  - harmonyos/entry/build.gradle"
echo "  - harmonyos/entry/package.json"
echo ""
echo "To build HarmonyOS package:"
echo "  1. Install OpenHarmony SDK"
echo "  2. Set OHOS_SDK_HOME environment variable"
echo "  3. Run: cd harmonyos && hvigorw build --mode module"
