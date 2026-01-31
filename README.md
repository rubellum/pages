# pages

Markdown を書くと HTML のサイトになる静的サイトジェネレーターです。**設定ファイルは不要**で、インストールしてすぐ使えます。

---

## インストール

```bash
go install github.com/rubellum/pages@latest
```

## はじめて使う（3ステップ）

**1. プロジェクトを用意する**

```bash
pages init
```

`src/`（Markdown 用）、`templates/`（HTML テンプレート）、`public/`（スタイル用）ができ、TOP ページのひな形 `src/index.md` が作成されます。

**2. ページを書く**

`src/` に `.md` ファイルを置きます。先頭にタイトルと日付を書く形式です。

```markdown
---
title: "はじめての記事"
date: 2025-01-15
---

ここに本文を書く。
```

**3. ビルドする**

```bash
pages build
```

`public/` に HTML が出力されます。このフォルダをそのまま Web サーバーで配信すればサイトの完成です。

---

## 日常の使い方

- **TOP ページ** … `src/index.md` の `title` がサイトのタイトルになります。
- **記事を増やす** … `src/` に新しい `.md` を追加するだけ。サブフォルダ（例: `src/blog/hello.md`）も使えます。
- **ビルド** … 更新のたびに `pages build` を実行。
- **プレビュー** … ビルド後、`public/` をローカルで開いて確認。例: `cd public && python -m http.server 3000` で http://localhost:3000 を開く。

---

## 設定は「あると便利」なだけ

設定ファイル `pages.yaml` は**なくても動きます**。サイトタイトルをファイルで変えたいときだけ、プロジェクトのルートに置きます。

```yaml
title: "わたしのブログ"
```

書かない項目はそのままデフォルトで動くので、必要なものだけ書けば大丈夫です。

---

## コマンド一覧

| コマンド | 説明 |
|----------|------|
| `pages init` | 新規プロジェクト用のフォルダとひな形を作成 |
| `pages build` | Markdown を HTML に変換（`src/` → `public/`） |

入出力フォルダを変えたいときは `pages build -i 入力フォルダ -o 出力フォルダ` で指定できます。

---

## ライセンス

MIT
