#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ITFeeds Docker 镜像构建脚本
功能：删除旧镜像，构建新镜像
"""

import subprocess
import sys
import time
from pathlib import Path

# 目录配置
ROOT = Path(__file__).resolve().parent.parent
IMAGE_NAME = "itfeeds"


def run(cmd, cwd=None):
    """运行命令并返回结果"""
    print(f"> {cmd}")
    try:
        result = subprocess.run(
            cmd,
            shell=True,
            cwd=cwd or ROOT,
            text=True
        )
        return result.returncode == 0
    except Exception as e:
        print(f"[FAIL] Command failed: {e}")
        return False


def check_docker():
    """检查 Docker 是否可用"""
    print("\n" + "=" * 50)
    print("  Check Docker")
    print("=" * 50)

    result = subprocess.run(
        "docker --version",
        shell=True,
        capture_output=True,
        text=True
    )

    if result.returncode != 0:
        print("[FAIL] Docker not available")
        print("Please ensure Docker is installed and running")
        return False

    print(f"[OK] {result.stdout.strip()}")
    return True


def remove_old_image():
    """删除旧镜像"""
    print("\n" + "=" * 50)
    print("  Remove Old Image")
    print("=" * 50)

    # 检查镜像是否存在
    result = subprocess.run(
        f"docker images -q {IMAGE_NAME}",
        shell=True,
        capture_output=True,
        text=True
    )

    image_ids = [img for img in result.stdout.strip().split('\n') if img]

    if not image_ids:
        print(f"[SKIP] Image '{IMAGE_NAME}' not found")
        return True

    # 删除镜像
    for img_id in image_ids:
        print(f"Removing image: {img_id}")
        if not run(f"docker rmi -f {img_id}"):
            print(f"[WARN] Failed to remove image: {img_id}")

    print("[OK] Old images removed")
    return True


def build_image():
    """构建新镜像"""
    print("\n" + "=" * 50)
    print("  Build Docker Image")
    print("=" * 50)

    if not run(f"docker build -t {IMAGE_NAME} ."):
        print("[FAIL] Docker build failed")
        return False

    print(f"[OK] Image '{IMAGE_NAME}' built successfully")
    return True


def show_image():
    """显示镜像信息"""
    print("\n" + "=" * 50)
    print("  Image Info")
    print("=" * 50 + "\n")

    run(f"docker images {IMAGE_NAME}")


def main():
    print("\n" + "=" * 50)
    print("  ITFeeds Docker Build Script")
    print("=" * 50)
    print(f"Image Name: {IMAGE_NAME}")
    print(f"Start Time: {time.strftime('%Y-%m-%d %H:%M:%S')}")

    # 1. 检查 Docker
    if not check_docker():
        return 1

    # 2. 删除旧镜像
    if not remove_old_image():
        return 1

    # 3. 构建新镜像
    if not build_image():
        return 1

    # 4. 显示镜像信息
    show_image()

    print("\n" + "=" * 50)
    print(f"End Time: {time.strftime('%Y-%m-%d %H:%M:%S')}")
    print("  Build Complete!")
    print("=" * 50 + "\n")

    return 0


if __name__ == "__main__":
    sys.exit(main())
