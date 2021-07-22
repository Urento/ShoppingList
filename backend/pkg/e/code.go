package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_EXIST_TAG       = 10001
	ERROR_EXIST_TAG_FAIL  = 10002
	ERROR_NOT_EXIST_TAG   = 10003
	ERROR_GET_TAGS_FAIL   = 10004
	ERROR_COUNT_TAG_FAIL  = 10005
	ERROR_ADD_TAG_FAIL    = 10006
	ERROR_EDIT_TAG_FAIL   = 10007
	ERROR_DELETE_TAG_FAIL = 10008
	ERROR_EXPORT_TAG_FAIL = 10009
	ERROR_IMPORT_TAG_FAIL = 10010

	ERROR_LIST_DOES_NOT_EXIST    = 10011
	ERROR_CHECK_EXIST_LIST_FAIL  = 10012
	ERROR_ADD_LIST_FAIL          = 10013
	ERROR_DELETE_LIST_FAIL       = 10014
	ERROR_EDIT_LIST_FAIL         = 10015
	ERROR_COUNT_LIST_FAIL        = 10016
	ERROR_GET_LISTS_FAIL         = 10017
	ERROR_GET_LIST_FAIL          = 10018
	ERROR_GETTING_LISTS_BY_OWNER = 10019

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
	ERROR_CREATING_ACCOUNT         = 20006

	ERROR_RATE_LIMITER        = 20007
	ERROR_RATELIMIT_TRY_LATER = 20008

	ERROR_GETTING_EMAIL_BY_JWT = 20009
	ERROR_RETRIEVING_USER_DATA = 20010
	ERROR_TOKEN_INVALID        = 20011

	ERROR_LIST_DOES_NOT_BELONG_TO_TOKEN = 20012
	ERROR_ID_IS_INVALID                 = 20013

	ERROR_NOT_AUTHORIZED = 20012
)
