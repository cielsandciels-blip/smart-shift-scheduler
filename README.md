# Smart Shift Scheduler (AI Shift Optimization System)

飲食店や小売店の「シフト作成業務」を自動化するSaaS型Webアプリケーションです。
店長が数時間かけて行っていたシフト組みのパズルを、数理最適化AI（Google OR-Tools）を用いて数秒で解決します。
「希望休」「連勤制限」「日別の必要人数」など、複雑な制約条件をすべて満たした最適解を自動生成するツールとして開発しました。

## 特徴 (Features)

1. **AI自動シフト生成**: Pythonの数理最適化ライブラリ（`Google OR-Tools`）を活用し、制約充足問題（CSP）として厳密な解を導出。最大5連勤までの労務コンプライアンスにも対応。
2. **直感的なUI/UX**: 生成されたシフトは `FullCalendar.js` 上でドラッグ＆ドロップするだけで修正可能。AI計算中は非同期でローディングを表示し、処理待ちのストレスを軽減。
3. **柔軟なルール設定**: 「土日は3人必要」「早番・遅番のバランス」など、店舗ごとに異なるルールをGUIから動的に追加・変更が可能。
4. **実務特化の出力**: 作成したシフトをExcelでそのまま開けるCSV形式（UTF-8 BOM付）で出力。月間の概算人件費もリアルタイムで試算。
5. **マイクロサービス構成**: 高速なGo言語サーバーと計算特化のPythonを連携させた、拡張性の高いアーキテクチャを採用。

## 🛠 使用技術 (Tech Stack)

- **Backend**: Go 1.23 / Gin (Clean Architecture採用)
- **AI Engine**: Python 3.9 / Google OR-Tools
- **Database**: PostgreSQL / SQLite (GORMにより切り替え可能)
- **Infrastructure**: Docker (Multi-stage build)
- **Frontend**: HTML5, CSS3, JavaScript (FullCalendar.js)
- **Deployment**: Render (Docker Container)

## プロジェクトの構造
- `backend/cmd/api/main.go`: Go APIサーバーのエントリーポイント
- `backend/internal/`: ビジネスロジック（Handler, Usecase, Repository）
- `backend/solver.py`: 数理最適化計算を行うPythonスクリプト
- `frontend/`: ユーザーインターフェース（HTML/CSS/JS）
- `Dockerfile`: GoとPythonが共存するコンテナ構築設定
- `go.mod`: Go言語の依存関係定義

## 開発の背景
これまで現場の店長が「経験と勘」に頼って長時間かけていた調整業務を、技術の力で効率化したいと考え開発しました。単にシフトを埋めるだけでなく、「過重労働の防止」や「人件費の管理」といった経営視点の課題も同時に解決できる、実用性を重視したソリューションを目指しました。

## 使い方 (Local Setup)
1. リポジトリをクローン
2. Dockerでデータベースを起動
   `docker run --name shift-db -p 5433:5432 -e POSTGRES_USER=manager -e POSTGRES_PASSWORD=manager123 -e POSTGRES_DB=shift_db -d postgres`
3. バックエンドディレクトリへ移動しサーバーを起動
   `cd backend && go run cmd/api/main.go`
4. ブラウザで `http://localhost:8080` にアクセス