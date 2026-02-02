from ortools.sat.python import cp_model
from .models import ShiftInput, ShiftResult

class ShiftSolver:
    def __init__(self, data: ShiftInput):
        self.data = data
        self.model = cp_model.CpModel()
        self.shifts = {} # (staff_id, day, shift_type) -> bool

    def solve(self) -> dict:
        # シフト区分: 0=休み, 1=早番, 2=遅番
        shift_types = [0, 1, 2] 
        num_days = self.data.days
        staff_ids = [s.id for s in self.data.staff_list]

        # 1. 変数の作成
        for s_id in staff_ids:
            for d in range(1, num_days + 1):
                # そのスタッフが、その日に「どのシフト(0,1,2)か」を表す変数
                self.shifts[(s_id, d)] = self.model.NewIntVar(0, 2, f'shift_s{s_id}_d{d}')

        # 2. 制約条件: 必要人数の確保
        # (今回は簡易的に「早番1人以上」「遅番1人以上」を毎日必須とする)
        for d in range(1, num_days + 1):
            early_shifts = []
            late_shifts = []
            for s_id in staff_ids:
                # shift_var == 1 (早番) かどうか判定する一時変数
                is_early = self.model.NewBoolVar(f'is_early_s{s_id}_d{d}')
                self.model.Add(self.shifts[(s_id, d)] == 1).OnlyEnforceIf(is_early)
                self.model.Add(self.shifts[(s_id, d)] != 1).OnlyEnforceIf(is_early.Not())
                early_shifts.append(is_early)

                # shift_var == 2 (遅番) かどうか判定する一時変数
                is_late = self.model.NewBoolVar(f'is_late_s{s_id}_d{d}')
                self.model.Add(self.shifts[(s_id, d)] == 2).OnlyEnforceIf(is_late)
                self.model.Add(self.shifts[(s_id, d)] != 2).OnlyEnforceIf(is_late.Not())
                late_shifts.append(is_late)

            # 毎日、早番1人以上、遅番1人以上
            self.model.Add(sum(early_shifts) >= 1)
            self.model.Add(sum(late_shifts) >= 1)

        # 3. 制約条件: 公平性 (勤務回数の平準化)
        # 全員の勤務回数の「最大値」と「最小値」の差を小さくする
        total_shifts_per_staff = []
        for s_id in staff_ids:
            # 0(休み)以外なら働くということ
            is_working_list = []
            for d in range(1, num_days + 1):
                is_working = self.model.NewBoolVar(f'working_s{s_id}_d{d}')
                self.model.Add(self.shifts[(s_id, d)] > 0).OnlyEnforceIf(is_working)
                self.model.Add(self.shifts[(s_id, d)] == 0).OnlyEnforceIf(is_working.Not())
                is_working_list.append(is_working)
            
            count_var = self.model.NewIntVar(0, num_days, f'count_{s_id}')
            self.model.Add(count_var == sum(is_working_list))
            total_shifts_per_staff.append(count_var)

        # 最大勤務数 - 最小勤務数 を最小化する
        max_shifts = self.model.NewIntVar(0, num_days, 'max_shifts')
        min_shifts = self.model.NewIntVar(0, num_days, 'min_shifts')
        self.model.AddMaxEquality(max_shifts, total_shifts_per_staff)
        self.model.AddMinEquality(min_shifts, total_shifts_per_staff)
        
        self.model.Minimize(max_shifts - min_shifts)

        # 4. 解く
        solver = cp_model.CpSolver()
        status = solver.Solve(self.model)

        # 5. 結果作成
        result_schedule = {}
        if status in [cp_model.OPTIMAL, cp_model.FEASIBLE]:
            for s_id in staff_ids:
                # 結果をリストで返す [1日目のシフト, 2日目のシフト...] (0:休み, 1:早, 2:遅)
                result_schedule[s_id] = [
                    solver.Value(self.shifts[(s_id, d)]) 
                    for d in range(1, num_days + 1)
                ]
            
            status_str = "Optimal" if status == cp_model.OPTIMAL else "Feasible"
            return {"status": status_str, "schedule": result_schedule}
        else:
            return {"status": "Infeasible", "schedule": {}}