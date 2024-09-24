package queryhelper

import (
	"fmt"

	model "github.com/fsidiqs/aegis-backend/model"
)

func SessionArrWhereProgramIDAndSeqNum(programIDSeqNum []model.ProgramIDSequenceNum) string {
	query := "SELECT * FROM m_program_sessions mps WHERE ("
	len := len(programIDSeqNum)
	for i := 0; i < len; i++ {
		if i == len-1 {
			query += fmt.Sprintf("(mps.program_id = '%v' AND mps.sequence_number = %v AND mps.record_flag = 'ACTIVE')", programIDSeqNum[i].ProgramID, programIDSeqNum[i].CurrentSequence)
		} else {
			query += fmt.Sprintf("(mps.program_id = '%v' AND mps.sequence_number = %v AND mps.record_flag = 'ACTIVE') OR ", programIDSeqNum[i].ProgramID, programIDSeqNum[i].CurrentSequence)
		}
	}
	query += ")"
	return query
}
