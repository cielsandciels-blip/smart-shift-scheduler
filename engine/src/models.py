from pydantic import BaseModel
from typing import List, Dict, Any

# --- 入力データの定義 ---
class Staff(BaseModel):
    id: int
    name: str
    is_leader: bool

class ShiftRequest(BaseModel):
    staff_id: int
    day: int
    priority: int

class DailyRequirement(BaseModel):
    day: int
    required_count: int
    required_leaders: int

class ShiftInput(BaseModel):
    staff_list: List[Staff]
    requests: List[ShiftRequest]
    requirements: List[DailyRequirement]
    days: int

# --- ★これを追加してください（出力データの定義） ---
class ShiftResult(BaseModel):
    status: str
    schedule: Dict[int, List[int]]