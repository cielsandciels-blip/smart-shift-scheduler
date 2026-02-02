import sys
import json
from src.models import ShiftInput
from src.solver import ShiftSolver
from pydantic import ValidationError

def main():
    try:
        # 1. 標準入力からJSONを受け取る (Go -> Python)
        input_str = sys.stdin.read()
        if not input_str:
            raise ValueError("Empty input")

        # 2. データのバリデーション (Pydantic)
        data = ShiftInput.model_validate_json(input_str)

        # 3. 計算実行
        solver = ShiftSolver(data)
        result = solver.solve()

        # 4. 結果をJSONとして標準出力へ (Python -> Go)
        print(json.dumps(result))

    except ValidationError as e:
        # 入力データ形式エラー
        error_res = {"status": "Error", "message": "Invalid Data Format", "details": e.errors()}
        print(json.dumps(error_res))
        sys.exit(1)
        
    except Exception as e:
        # その他のシステムエラー
        error_res = {"status": "Error", "message": str(e)}
        print(json.dumps(error_res))
        sys.exit(1)

if __name__ == "__main__":
    main()