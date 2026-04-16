#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ITFeeds 二进制构建脚本
功能：
1. gf pack - 打包静态资源到 packed.go
2. gf build - 构建二进制可执行文件
"""

import subprocess
import sys
import time
from pathlib import Path

# 目录配置
ROOT = Path(__file__).resolve().parent.parent
RESOURCE = ROOT / "resource"
OUTPUT_DIR = ROOT / "output"


def run(cmd, cwd=None):
    """运行命令"""
    print(f"\n> {cmd}")
    print("-" * 50)
    result = subprocess.run(cmd, shell=True, cwd=cwd or ROOT)
    return result.returncode == 0


def pack_resource():
    """gf pack - 打包静态资源"""
    print("\n" + "=" * 50)
    print("  Step 1: Pack Resource")
    print("=" * 50)

    # 检查 gf 命令
    if not run("gf version"):
        print("[FAIL] gf command not found")
        print("Please install: go install github.com/gogf/gf/cmd/gf/v2@latest")
        return False

    # 打包 resource 目录
    cmd = f'gf pack "{RESOURCE}" internal/packed/packed.go -n packed'
    if not run(cmd):
        print("[FAIL] gf pack failed")
        return False

    print("[OK] Resource packed to internal/packed/packed.go")
    return True


def build_binary():
    """gf build - 构建二进制"""
    print("\n" + "=" * 50)
    print("  Step 2: Build Binary")
    print("=" * 50)

    # 创建输出目录
    OUTPUT_DIR.mkdir(exist_ok=True)

    # 构建二进制（Linux AMD64）
    cmd = "gf build main.go -n itfeeds -a amd64 -s linux -p output"
    if not run(cmd):
        print("[FAIL] gf build failed")
        return False

    print("[OK] Binary built to output/itfeeds")
    return True


def show_result():
    """显示构建结果"""
    print("\n" + "=" * 50)
    print("  Build Result")
    print("=" * 50 + "\n")

    # 显示二进制文件信息
    binary = OUTPUT_DIR / "itfeeds"
    if binary.exists():
        print(f"[OK] Binary: {binary}")
        run(f"ls -lh {binary}")
    else:
        print("[FAIL] Binary not found")

    # 显示 packed 文件
    packed = ROOT / "internal" / "packed" / "packed.go"
    if packed.exists():
        size = packed.stat().st_size / 1024  # KB
        print(f"\n[OK] Packed: {packed} ({size:.1f} KB)")


def main():
    print("\n" + "=" * 50)
    print("  ITFeeds Binary Build Script")
    print("=" * 50)
    print(f"开始时间: {time.strftime('%Y-%m-%d %H:%M:%S')}")

    # 1. 打包静态资源
    if not pack_resource():
        return 1

    # 2. 构建二进制
    if not build_binary():
        return 1

    # 3. 显示结果
    show_result()

    print("\n" + "=" * 50)
    print(f"完成时间: {time.strftime('%Y-%m-%d %H:%M:%S')}")
    print("  Build Complete!")
    print("=" * 50 + "\n")

    return 0


if __name__ == "__main__":
    sys.exit(main())
