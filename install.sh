#!/bin/bash

set -e

clear
cat <<"EOF"
                  ////////  
             /////  ////////  
           //////// ////////  
 //////// ///////// ///////   
 //////// ///////// //////    
 //////// ///////// ////      
 //////// ///////// ///       
 //////// ///////// /////     
 //////// ///////// //////    
 //////// ///////// ////////  
 //////// ///////// ////////  
 //////// ///////// ////////  
 //////// ////////            
 ////////  /////              
 ///////                      

All you need. Nothing you don't.

EOF

OS_NAME=$(uname -s)
CPU_ARCH=$(uname -m)
SILENT=false
AUTO_SETUP=false
DOMAIN_NAME=""

while [ "$#" -gt 0 ]; do
  case "$1" in
    --silent)
      SILENT=true
      shift 1
      ;;
    --setup)
      AUTO_SETUP=true
      shift 1
      ;;
    --domain)
      DOMAIN_NAME="$2"
      shift 2
      ;;
    --domain=*)
      DOMAIN_NAME="${1#*=}"
      shift 1
      ;;
    *)
      echo "Unknown option: $1" >&2
      exit 1
      ;;
  esac
done

if [ "${OS_NAME}" != "Linux" ] && [ "${OS_NAME}" != "Darwin" ]; then
    echo "drim only works on Linux and macOS"
    exit 1
fi

get_binary_name() {
    local platform=""
    local arch=""

    case "${OS_NAME}" in
        Linux)
            platform="Linux"
            ;;
        Darwin)
            platform="Darwin"
            ;;
    esac

    case "${CPU_ARCH}" in
        x86_64|amd64)
            arch="x86_64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            echo "Unsupported architecture: ${CPU_ARCH}"
            exit 1
            ;;
    esac

    echo "drim_${platform}_${arch}"
}

BINARY_NAME=$(get_binary_name)
DOWNLOAD_URL="https://github.com/usekaneo/drim/releases/latest/download/${BINARY_NAME}"

if [ "${SILENT}" = false ]; then
    echo "Downloading drim for ${OS_NAME} ${CPU_ARCH}..."
fi

TEMP_FILE=$(mktemp)
trap "rm -f ${TEMP_FILE}" EXIT

HTTP_CODE=$(curl -sL "${DOWNLOAD_URL}" -o "${TEMP_FILE}" -w "%{http_code}")

if [ "${HTTP_CODE}" != "200" ]; then
    echo "Download failed with HTTP code: ${HTTP_CODE}"
    echo "URL: ${DOWNLOAD_URL}"
    exit 1
fi

chmod +x "${TEMP_FILE}"

INSTALL_DIR="/usr/local/bin"
if [ "${OS_NAME}" = "Linux" ]; then
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    else
        INSTALL_DIR="/bin"
    fi
fi

if [ "${SILENT}" = false ]; then
    echo "Installing drim to ${INSTALL_DIR}..."
fi

if [ -w "${INSTALL_DIR}" ]; then
    mv "${TEMP_FILE}" "${INSTALL_DIR}/drim"
else
    sudo mv "${TEMP_FILE}" "${INSTALL_DIR}/drim"
fi

if [ "${SILENT}" = false ]; then
    echo ""
    echo "✅ drim installed successfully!"
    echo ""
    drim --version
    echo ""
fi

if [ "${AUTO_SETUP}" = true ]; then
    if [ "${SILENT}" = false ]; then
        echo "Running setup..."
        echo ""
    fi
    
    if [ -n "${DOMAIN_NAME}" ]; then
        echo "${DOMAIN_NAME}" | drim setup
    else
        drim setup
    fi
else
    if [ "${SILENT}" = false ]; then
        echo "Run 'drim setup' to deploy Kaneo"
    fi
fi

