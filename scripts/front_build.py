#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ITFeeds 前端构建脚本
功能：构建 Vue 前端并输出到 resource/public
"""

import shutil
import subprocess
import sys
import time
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
WEB_DIR = ROOT / "web"
PUBLIC_DIR = ROOT / "resource" / "public"


def run(cmd, cwd=None):
    """运行命令"""
    print(f"\n> {cmd}")
    print("-" * 50)
    result = subprocess.run(cmd, shell=True, cwd=cwd or ROOT)
    return result.returncode == 0


def clean_public():
    """清空 resource/public"""
    print("\n" + "=" * 50)
    print("  Step 1: Clean resource/public")
    print("=" * 50)

    if PUBLIC_DIR.exists():
        shutil.rmtree(PUBLIC_DIR)
        print(f"[OK] Removed {PUBLIC_DIR}")
    else:
        print("[OK] Directory does not exist, skip")

    PUBLIC_DIR.mkdir(parents=True, exist_ok=True)
    return True


def build_frontend():
    """构建前端"""
    print("\n" + "=" * 50)
    print("  Step 2: Build Frontend")
    print("=" * 50)

    if not run("npm install", cwd=WEB_DIR):
        print("[FAIL] npm install failed")
        return False

    if not run("npm run build", cwd=WEB_DIR):
        print("[FAIL] npm run build failed")
        return False

    return True


def copy_to_public():
    """将构建产物复制到 resource/public"""
    print("\n" + "=" * 50)
    print("  Step 3: Copy to resource/public")
    print("=" * 50)

    dist_dir = WEB_DIR / "dist"
    if not dist_dir.exists():
        print(f"[FAIL] {dist_dir} not found")
        return False

    for item in dist_dir.iterdir():
        dest = PUBLIC_DIR / item.name
        if item.is_dir():
            shutil.copytree(item, dest)
        else:
            shutil.copy2(item, dest)
        print(f"  {item.name}")

    print(f"\n[OK] Copied to {PUBLIC_DIR}")
    return True


def main():
    print("\n" + "=" * 50)
    print("  ITFeeds Frontend Build")
    print("=" * 50)
    print(f"Time: {time.strftime('%Y-%m-%d %H:%M:%S')}")

    if not clean_public():
        return 1
    if not build_frontend():
        return 1
    if not copy_to_public():
        return 1

    print("\n" + "=" * 50)
    print("  Build Complete!")
    print("=" * 50 + "\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
