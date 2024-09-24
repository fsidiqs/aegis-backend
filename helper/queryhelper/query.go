package queryhelper

const (
	ActiveFlag         = "record_flag = 'ACTIVE'"
	ActiveAndPublished = "record_flag = 'ACTIVE' AND status = 'PUBLISHED'"
)

// Program
const (
	//-- m_programs
	SelectProgramByCategoryLimit1 = "SELECT * FROM m_programs WHERE category = ? AND status = 'PUBLISHED' AND record_flag = 'ACTIVE' LIMIT 1"
	// param: user_id
	ProgramsFromEnrollments = `SELECT peh.id AS "peh.id", peh.program_id AS "peh.program_id", peh.user_id AS "peh.user_id", peh.status AS "peh.status", peh.current_sequence AS "peh.current_sequence", peh.rating AS "peh.rating", peh.review AS "peh.review", peh.created_at AS "peh.created_at", peh.updated_at AS "peh.updated_at", peh.created_by AS "peh.created_by", peh.updated_by AS "peh.updated_by", peh.record_flag AS "peh.record_flag", mp.id AS "mp.id", mp."name" AS "mp.name", mp.category AS "mp.category", mp.sub_category AS "mp.sub_category", mp.description AS "mp.description", mp.objective AS "mp.objective", mp.thumbnail_url AS "mp.thumbnail_url", mp.rating AS "mp.rating", mp.duration AS "mp.duration", mp.subscription_type AS "mp.subscription_type", mp.intro_video_id AS "mp.intro_video_id", mp.created_at AS "mp.created_at", mp.updated_at AS "mp.updated_at", mp.created_by AS "mp.created_by", mp.updated_by AS "mp.updated_by", mp.record_flag AS "mp.record_flag", mp.coach_id AS "mp.coach_id", mp.status AS "mp.status", mp.content_type AS "mp.content_type" FROM program_enrollment_history peh INNER JOIN m_programs mp ON peh.program_id = mp.id WHERE peh.user_id = ? AND peh.record_flag = 'ACTIVE' AND peh.status = 'IN_PROGRESS' AND mp.record_flag = 'ACTIVE' AND mp.status = 'PUBLISHED'`
	// m_program_videos
	SelectProgramVideoWithAssets = "SELECT * FROM m_program_videos mpv where mpv.id = ?"
	// program_enrollment_history
	CountUserFreeEnrollment                             = "SELECT count(peh.user_id) FROM program_enrollment_history peh INNER JOIN m_programs mp ON peh.program_id  = mp.id WHERE peh.user_id = '?' AND mp.subscription_type = 'FREE'"
	IsUserEnrollmentInProgressToProgamID                = "SELECT EXISTS(SELECT 1 FROM program_enrollment_history WHERE user_id = ? AND program_id = ? AND record_flag = 'ACTIVE' AND status = 'IN_PROGRESS' LIMIT 1)"
	GetEnrollmentInProgressByUserIDProgramID            = "SELECT * FROM program_enrollment_history peh WHERE peh.user_id = ? AND program_id = ? AND status = 'IN_PROGRESS' AND record_flag = 'ACTIVE'"
	GetEnrollmentInProgressByUserIDProgramIDLoadProgram = `SELECT peh.id AS "peh.id", peh.program_id AS "peh.program_id", peh.user_id AS "peh.user_id", peh.status AS "peh.status", peh.current_sequence AS "peh.current_sequence", peh.rating AS "peh.rating", peh.review AS "peh.review", peh.created_at AS "peh.created_at", peh.updated_at AS "peh.updated_at", peh.created_by AS "peh.created_by", peh.updated_by AS "peh.updated_by", peh.record_flag AS "peh.record_flag", mp.id AS "mp.id", mp."name" AS "mp.name", mp.category AS "mp.category", mp.sub_category AS "mp.sub_category", mp.description AS "mp.description", mp.objective AS "mp.objective", mp.thumbnail_url AS "mp.thumbnail_url", mp.rating AS "mp.rating", mp.duration AS "mp.duration", mp.subscription_type AS "mp.subscription_type", mp.intro_video_id AS "mp.intro_video_id", mp.created_at AS "mp.created_at", mp.updated_at AS "mp.updated_at", mp.created_by AS "mp.created_by", mp.updated_by AS "mp.updated_by", mp.record_flag AS "mp.record_flag", mp.coach_id AS "mp.coach_id", mp.status AS "mp.status" FROM program_enrollment_history peh INNER JOIN m_programs mp ON  peh.program_id = mp.id WHERE peh.user_id = ? AND peh.program_id = ? AND peh.status = 'IN_PROGRESS' AND peh.record_flag = 'ACTIVE' AND mp.status = 'PUBLISHED' AND mp.record_flag = 'ACTIVE'`
	GetUserEnrollmentInProgressLoadProgram              = `SELECT peh.id AS "peh.id", peh.program_id AS "peh.program_id", peh.user_id AS "peh.user_id", peh.status AS "peh.status", peh.current_sequence AS "peh.current_sequence", peh.rating AS "peh.rating", peh.review AS "peh.review", peh.created_at AS "peh.created_at", peh.updated_at AS "peh.updated_at", peh.created_by AS "peh.created_by", peh.updated_by AS "peh.updated_by", peh.record_flag AS "peh.record_flag", mp.id AS "mp.id", mp."name" AS "mp.name", mp.category AS "mp.category", mp.sub_category AS "mp.sub_category", mp.description AS "mp.description", mp.objective AS "mp.objective", mp.thumbnail_url AS "mp.thumbnail_url", mp.rating AS "mp.rating", mp.duration AS "mp.duration", mp.subscription_type AS "mp.subscription_type", mp.intro_video_id AS "mp.intro_video_id", mp.created_at AS "mp.created_at", mp.updated_at AS "mp.updated_at", mp.created_by AS "mp.created_by", mp.updated_by AS "mp.updated_by", mp.record_flag AS "mp.record_flag", mp.coach_id AS "mp.coach_id", mp.status AS "mp.status" FROM program_enrollment_history peh INNER JOIN m_programs mp ON peh.program_id = mp.id WHERE peh.user_id = ? AND peh.record_flag = 'ACTIVE' AND peh.status = 'IN_PROGRESS' AND mp.record_flag = 'ACTIVE' AND mp.status = 'PUBLISHED'`
	// m_program_sessions
	// get the last sequence_number of a program
	LastSessionSeqNum = "SELECT mps.sequence_number FROM m_program_sessions mps WHERE mps.program_id = ? ORDER BY mps.sequence_number DESC LIMIT 1"

	// m_program_activities
	GetProgramActivitiesPluckIDByProgramSessionID = "SELECT mpa.id FROM m_program_sessions mps INNER JOIN m_program_sections mpsect ON mps.id = mpsect.program_session_id INNER JOIN m_program_activities mpa ON mpsect.id = mpa.program_section_id WHERE mps.id = ? AND mps.record_flag = 'ACTIVE' AND mps.record_flag = 'ACTIVE' AND mpsect.record_flag = 'ACTIVE' AND mpa.record_flag = 'ACTIVE'"
)

const (
	EmailExists       = "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? LIMIT 1)"
	SocialLoginExists = "SELECT EXISTS(SELECT 1 FROM social_logins WHERE social_user_id = ? AND provider = ? LIMIT 1)"
)

// PAYMENT
const (
	UserHasUnfinishedPayment = "SELECT EXISTS(SELECT 1 FROM payments WHERE user_id = ? AND status = 'WAITING_PAYMENT' AND record_flag = 'ACTIVE' LIMIT 1)"
)

// SUBSCRIPTION
const (
	IsUserSubscribing = "SELECT EXISTS(SELECT 1 FROM subscriptions WHERE user_id = ? AND end_date >= ? AND record_flag = 'ACTIVE' ORDER BY created_at DESC LIMIT 1)"
)

// PROMO
const (
	GetRefAcquisitionByDeviceIDRefCode = `SELECT r.id AS "r.id", r.referral_campaign_id AS "r.referral_campaign_id", r.referrer_id AS "r.referrer_id", r.referral_code AS "r.referral_code", r.created_at AS "r.created_at", r.updated_at AS "r.updated_at", r.created_by AS "r.created_by", r.updated_by AS "r.updated_by", r.record_flag AS "r.record_flag", ra.id AS "ra.id", ra.referral_id AS "ra.referral_id", ra.referee_id AS "ra.referee_id", ra.device_id AS "ra.device_id", ra.created_at AS "ra.created_at", ra.updated_at AS "ra.updated_at", ra.created_by AS "ra.created_by", ra.updated_by AS "ra.updated_by", ra.record_flag AS "ra.record_flag" FROM referrals r INNER JOIN referral_acquisitions ra ON r.id = ra.referral_id WHERE r.referral_code = ? AND r.record_flag = ? AND ra.device_id = ? AND ra.record_flag = ?`
	IsDeviceIDAndReferralCodeExist     = "SELECT EXISTS(SELECT 1 FROM referrals r INNER JOIN referral_acquisitions ra ON r.id = ra.referral_id WHERE ra.device_id = ? AND ra.record_flag = 'ACTIVE' AND r.referral_code = ? AND r.record_flag = 'ACTIVE' LIMIT 1)"
	IsReferralCodeExists               = "SELECT EXISTS(SELECT 1 FROM referrals r WHERE r.referral_code = ? AND r.record_flag = 'ACTIVE' LIMIT 1)"
)
