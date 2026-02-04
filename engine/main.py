import sys
import json
from ortools.sat.python import cp_model

def solve_shift_scheduling(input_data):
    # データの読み込み
    staff_list = input_data.get('staff_list', [])
    days = input_data.get('days', 30)
    requests = input_data.get('requests', []) 
    
    # ★ここを修正しました ( or [] を追加 )
    role_constraints = input_data.get('role_constraints') or []

    # モデル作成
    model = cp_model.CpModel()
    shifts = {} 
    shift_types = [0, 1, 2]

    # 1. 変数の定義
    for s in staff_list:
        s_id = s['id']
        for d in range(days):
            for st in shift_types:
                shifts[(s_id, d, st)] = model.NewBoolVar(f'shift_s{s_id}_d{d}_st{st}')
            model.Add(sum(shifts[(s_id, d, st)] for st in shift_types) == 1)

    # 2. 希望休 (NG) の反映 (簡易実装: 今回はスキップ)

    # 3. 基本人数要件
    for d in range(days):
        model.Add(sum(shifts[(s['id'], d, 1)] for s in staff_list) >= 1)
        model.Add(sum(shifts[(s['id'], d, 2)] for s in staff_list) >= 1)

    # 4. 役割（スキル）の制約
    for constraint in role_constraints:
        target_role = constraint.get('role', '') 
        required_count = constraint.get('count', 0) 

        if not target_role or required_count <= 0:
            continue

        qualified_staff_ids = []
        for s in staff_list:
            roles = s.get('roles', '')
            if target_role in roles: 
                qualified_staff_ids.append(s['id'])
        
        if len(qualified_staff_ids) > 0:
            for d in range(days):
                model.Add(sum(shifts[(sid, d, 1)] for sid in qualified_staff_ids) >= required_count)
                model.Add(sum(shifts[(sid, d, 2)] for sid in qualified_staff_ids) >= required_count)

    # 5. 目的関数
    total_shifts = sum(shifts[(s['id'], d, st)] 
                       for s in staff_list for d in range(days) for st in [1, 2])
    model.Maximize(total_shifts)

    # 解く
    solver = cp_model.CpSolver()
    status = solver.Solve(model)

    result = {"schedule": {}, "status": solver.StatusName(status)}
    
    if status == cp_model.OPTIMAL or status == cp_model.FEASIBLE:
        for s in staff_list:
            s_id = s['id']
            daily_shifts = []
            for d in range(days):
                for st in shift_types:
                    if solver.Value(shifts[(s_id, d, st)]) == 1:
                        daily_shifts.append(st)
                        break
            result["schedule"][s_id] = daily_shifts

    return result

if __name__ == '__main__':
    input_json = sys.stdin.read()
    data = json.loads(input_json)
    result = solve_shift_scheduling(data)
    print(json.dumps(result))