#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ITFeeds Docker 镜像导出脚本
功能：将镜像导出为 tar 文件
"""

import subprocess
import sys
import time
import argparse
from pathlib import Path
from datetime import datetime

# 默认配置
IMAGE_NAME = "itfeeds"
ROOT = Path(__file__).resolve().parent.parent


def run(cmd, cwd=None):
    """运行命令"""
    print(f"> {cmd}")
    result = subprocess.run(cmd, shell=True, cwd=cwd or ROOT)
    return result.returncode == 0


def check_image_exists(image_name):
    """检查镜像是否存在"""
    result = subprocess.run(
        f"docker images -q {image_name}",
        shell=True,
        capture_output=True,
        text=True
    )
    return bool(result.stdout.strip())


def export_image(image_name, output_path):
    """导出镜像"""
    print("\n" + "=" * 50)
    print("  Export Docker Image")
    print("=" * 50)
    print(f"Image: {image_name}")
    print(f"Output: {output_path}")
    print("=" * 50 + "\n")

    # 检查镜像
    if not check_image_exists(image_name):
        print(f"[FAIL] Image '{image_name}' not found")
        print("Please run 'python build_docker.py' first")
        return False

    # 导出镜像
    cmd = f"docker save -o \"{output_path}\" {image_name}"
    if not run(cmd):
        print("[FAIL] Export failed")
        return False

    # 显示文件信息
    if Path(output_path).exists():
        size = Path(output_path).stat().st_size / (1024 * 1024)  # MB
        print(f"\n[OK] Exported: {output_path}")
        print(f"[OK] Size: {size:.1f} MB")
        return True
    else:
        print("[FAIL] Output file not created")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="导出 Docker 镜像为 tar 文件",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  python export_image.py                    # 导出到当前目录
  python export_image.py -o /tmp/itfeeds.tar # 导出到指定路径
  python export_image.py -n myimage         # 导出指定镜像
        """
    )
    parser.add_argument(
        "-o", "--output",
        default=None,
        help="输出路径 (默认: ./itfeeds_YYYYMMDD_HHMMSS.tar)"
    )
    parser.add_argument(
        "-n", "--name",
        default=IMAGE_NAME,
        help=f"镜像名称 (默认: {IMAGE_NAME})"
    )

    args = parser.parse_args()

    # 确定输出路径
    if args.output:
        output_path = Path(args.output)
    else:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output_path = ROOT / f"{args.name}_{timestamp}.tar"

    # 确保输出目录存在
    output_path.parent.mkdir(parents=True, exist_ok=True)

    print("\n" + "=" * 50)
    print("  ITFeeds Docker Image Export")
    print("=" * 50)
    print(f"Time: {time.strftime('%Y-%m-%d %H:%M:%S')}")

    # 导出镜像
    if not export_image(args.name, str(output_path)):
        return 1

    print("\n" + "=" * 50)
    print("  Export Complete!")
    print("=" * 50 + "\n")

    return 0


if __name__ == "__main__":
    sys.exit(main())
