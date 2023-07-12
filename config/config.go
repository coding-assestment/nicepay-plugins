package config

const TESTING_REQ_URL = "https://staging.nicepay.co.id"
const NICEPAY_REQ_URL = "https://staging.nicepay.co.id/nicepay/direct/v2/registration"
const NICEPAY_REQ_CC_URL = "https://staging.nicepay.co.id/nicepay/direct/v2/registration"
const NICEPAY_CANCEL_URL = "https://staging.nicepay.co.id/nicepay/direct/v2/cancel"
const NICEPAY_ORDER_STATUS_URL = "https://staging.nicepay.co.id/nicepay/direct/v2/inquiry"

const NICEPAY_IMID = "IONPAYTEST"                                                                                       // Merchant ID
const NICEPAY_MERCHANT_KEY = "33F49GnCMS1mFYlGXisbUDzVf2ATWCl9k3R++d5hDd3Frmuos/XLx8XhXpe+LDYAbpGKZYSwtlyyLOtS/8aD7A==" // API Key Merchant Key

// const NICEPAY_CALLBACK_URL =       "http://".$_SERVER['SERVER_NAME'].dirname($_SERVER['PHP_SELF'])."/resultGate.php");      // Merchant's result page URL
const NICEPAY_DBPROCESS_URL = "http://httpresponder.com/nicepay" // Merchant's notification handler URL

/* TIMEOUT - Define as needed (in seconds) */
const NICEPAY_TIMEOUT_CONNECT = 15
const NICEPAY_TIMEOUT_READ = 25

const NICEPAY_LOG_CRITICAL = 1
const NICEPAY_LOG_ERROR = 2
const NICEPAY_LOG_NOTICE = 3
const NICEPAY_LOG_INFO = 5
const NICEPAY_LOG_DEBUG = 7
