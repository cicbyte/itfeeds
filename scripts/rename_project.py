#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
项目重命名脚本
将 itfeeds 重命名为 itfeeds，修正 cicbyte → cicbyte

替换规则：
  1. github.com/cicbyte/itfeeds → github.com/cicbyte/itfeeds  (Go import 路径)
  2. itfeeds → itfeeds                                        (二进制名/镜像名)
  3. cicbyte → cicbyte                                          (GitHub 用户名修正)

跳过规则：
  - node_modules/
  - vendor/
  - .git/
  - *.exe / *.sum (go.sum 会由 go mod tidy 重新生成)
"""

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent

# 需要跳过的目录
SKIP_DIRS = {"node_modules", "vendor", ".git", ".idea", ".vscode", "output"}

# 文件扩展名白名单（只处理文本文件）
TEXT_EXTENSIONS = {
    ".go", ".mod", ".sum", ".yaml", ".yml", ".json", ".py",
    ".md", ".txt", ".sql", ".toml", ".xml", ".sh", ".bat",
    ".html", ".css", ".js", ".ts", ".vue", ".jsx", ".tsx",
    ".env", ".gitignore", ".dockerignore",
}

# 无扩展名的文件也处理（如 Dockerfile, Makefile）
TEXT_NAMES = {"Dockerfile", "Makefile", ".gitignore", ".dockerignore", "CLAUDE.md"}


def should_process(path: Path) -> bool:
    """判断文件是否需要处理"""
    # 检查是否在跳过目录中
    for part in path.parts:
        if part in SKIP_DIRS:
            return False

    # 检查是否是文本文件
    if path.name in TEXT_NAMES:
        return True
    if path.suffix.lower() in TEXT_EXTENSIONS:
        return True
    return False


def replace_in_file(filepath: Path, replacements: list[tuple[str, str]]) -> tuple[int, list[str]]:
    """对文件执行替换，返回 (修改次数, 详情列表)"""
    try:
        content = filepath.read_text(encoding="utf-8")
    except (UnicodeDecodeError, PermissionError):
        return 0, []

    original = content
    changes = []

    for old, new in replacements:
        if old in content:
            count = content.count(old)
            content = content.replace(old, new)
            changes.append(f"  {old} → {new} ({count} 处)")

    if content != original:
        filepath.write_text(content, encoding="utf-8")
        return len(changes), changes

    return 0, []


def main():
    # 定义替换规则（顺序很重要，先替换长的再替换短的）
    replacements = [
        ("github.com/cicbyte/itfeeds", "github.com/cicbyte/itfeeds"),
        ("itfeeds", "itfeeds"),
    ]

    # cicbyte 单独处理（避免影响已经替换过的 github.com/cicbyte/itfeeds）
    replacements_no_cicbyte = replacements[:]
    replacements_with_cicbyte = [
        ("cicbyte", "cicbyte"),
    ] + replacements

    total_files = 0
    total_changes = 0

    print("=" * 60)
    print("  项目重命名: itfeeds → itfeeds")
    print("=" * 60)
    print()

    # 收集所有文件
    all_files = sorted(ROOT.rglob("*"))

    for filepath in all_files:
        if not filepath.is_file():
            continue
        if not should_process(filepath):
            continue

        # go.sum 由 go mod tidy 重新生成，跳过
        if filepath.name == "go.sum":
            continue

        rel = filepath.relative_to(ROOT)

        # 先检查是否包含 cicbyte（需要在 import 路径替换之前处理）
        try:
            content = filepath.read_text(encoding="utf-8")
        except (UnicodeDecodeError, PermissionError):
            continue

        if "cicbyte" in content:
            n, details = replace_in_file(filepath, replacements_with_cicbyte)
        else:
            n, details = replace_in_file(filepath, replacements_no_cicbyte)

        if n > 0:
            total_files += 1
            total_changes += n
            print(f"[MODIFIED] {rel}")
            for d in details:
                print(d)
            print()

    print("=" * 60)
    print(f"  完成: 修改了 {total_files} 个文件, {total_changes} 处替换")
    print("=" * 60)
    print()
    print("后续步骤:")
    print("  1. go mod tidy          # 重新生成 go.sum")
    print("  2. gf gen dao           # 重新生成 DAO（如有需要）")
    print("  3. 验证构建: go build ./...")


if __name__ == "__main__":
    main()
