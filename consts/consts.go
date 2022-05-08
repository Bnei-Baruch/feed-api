package consts

const (
	// Collection Types
	CT_ARTICLES           = "ARTICLES"
	CT_BOOKS              = "BOOKS"
	CT_CHILDREN_LESSONS   = "CHILDREN_LESSONS"
	CT_CLIPS              = "CLIPS"
	CT_CONGRESS           = "CONGRESS"
	CT_DAILY_LESSON       = "DAILY_LESSON"
	CT_FRIENDS_GATHERINGS = "FRIENDS_GATHERINGS"
	CT_HOLIDAY            = "HOLIDAY"
	CT_LECTURE_SERIES     = "LECTURE_SERIES"
	CT_LESSONS_SERIES     = "LESSONS_SERIES"
	CT_MEALS              = "MEALS"
	CT_PICNIC             = "PICNIC"
	CT_SONGS              = "SONGS"
	CT_SPECIAL_LESSON     = "SPECIAL_LESSON"
	CT_UNITY_DAY          = "UNITY_DAY"
	CT_VIDEO_PROGRAM      = "VIDEO_PROGRAM"
	CT_VIRTUAL_LESSONS    = "VIRTUAL_LESSONS"
	CT_WOMEN_LESSONS      = "WOMEN_LESSONS"

	// Content Unit Types
	CT_ARTICLE               = "ARTICLE"
	CT_BLOG_POST             = "BLOG_POST"
	CT_BOOK                  = "BOOK"
	CT_CHILDREN_LESSON       = "CHILDREN_LESSON"
	CT_CLIP                  = "CLIP"
	CT_EVENT_PART            = "EVENT_PART"
	CT_FRIENDS_GATHERING     = "FRIENDS_GATHERING"
	CT_FULL_LESSON           = "FULL_LESSON"
	CT_KITEI_MAKOR           = "KITEI_MAKOR"
	CT_LECTURE               = "LECTURE"
	CT_LELO_MIKUD            = "LELO_MIKUD"
	CT_LESSON_PART           = "LESSON_PART"
	CT_MEAL                  = "MEAL"
	CT_PUBLICATION           = "PUBLICATION"
	CT_RESEARCH_MATERIAL     = "RESEARCH_MATERIAL"
	CT_SONG                  = "SONG"
	CT_TRAINING              = "TRAINING"
	CT_UNKNOWN               = "UNKNOWN"
	CT_VIDEO_PROGRAM_CHAPTER = "VIDEO_PROGRAM_CHAPTER"
	CT_VIRTUAL_LESSON        = "VIRTUAL_LESSON"
	CT_WOMEN_LESSON          = "WOMEN_LESSON"

	// Security levels
	SEC_PUBLIC    = int16(0)
	SEC_SENSITIVE = int16(1)
	SEC_PRIVATE   = int16(2)

	// Languages
	LANG_ENGLISH    = "en"
	LANG_HEBREW     = "he"
	LANG_RUSSIAN    = "ru"
	LANG_SPANISH    = "es"
	LANG_ITALIAN    = "it"
	LANG_GERMAN     = "de"
	LANG_DUTCH      = "nl"
	LANG_FRENCH     = "fr"
	LANG_PORTUGUESE = "pt"
	LANG_TURKISH    = "tr"
	LANG_POLISH     = "pl"
	LANG_ARABIC     = "ar"
	LANG_HUNGARIAN  = "hu"
	LANG_FINNISH    = "fi"
	LANG_LITHUANIAN = "lt"
	LANG_JAPANESE   = "ja"
	LANG_BULGARIAN  = "bg"
	LANG_GEORGIAN   = "ka"
	LANG_NORWEGIAN  = "no"
	LANG_SWEDISH    = "sv"
	LANG_CROATIAN   = "hr"
	LANG_CHINESE    = "zh"
	LANG_PERSIAN    = "fa"
	LANG_ROMANIAN   = "ro"
	LANG_HINDI      = "hi"
	LANG_UKRAINIAN  = "ua"
	LANG_MACEDONIAN = "mk"
	LANG_SLOVENIAN  = "sl"
	LANG_LATVIAN    = "lv"
	LANG_SLOVAK     = "sk"
	LANG_CZECH      = "cs"
	LANG_AMHARIC    = "am"
	LANG_MULTI      = "zz"
	LANG_UNKNOWN    = "xx"
)

var I18N_LANG_ORDER = map[string][]string{
	"":              {LANG_ENGLISH},
	LANG_ENGLISH:    {LANG_ENGLISH},
	LANG_HEBREW:     {LANG_HEBREW, LANG_ENGLISH},
	LANG_RUSSIAN:    {LANG_RUSSIAN, LANG_ENGLISH},
	LANG_SPANISH:    {LANG_SPANISH, LANG_ENGLISH},
	LANG_ITALIAN:    {LANG_ITALIAN, LANG_ENGLISH},
	LANG_GERMAN:     {LANG_GERMAN, LANG_ENGLISH},
	LANG_DUTCH:      {LANG_DUTCH, LANG_ENGLISH},
	LANG_FRENCH:     {LANG_FRENCH, LANG_ENGLISH},
	LANG_PORTUGUESE: {LANG_PORTUGUESE, LANG_ENGLISH},
	LANG_TURKISH:    {LANG_TURKISH, LANG_ENGLISH},
	LANG_POLISH:     {LANG_POLISH, LANG_ENGLISH},
	LANG_ARABIC:     {LANG_ARABIC, LANG_ENGLISH},
	LANG_HUNGARIAN:  {LANG_HUNGARIAN, LANG_ENGLISH},
	LANG_FINNISH:    {LANG_FINNISH, LANG_ENGLISH},
	LANG_LITHUANIAN: {LANG_LITHUANIAN, LANG_RUSSIAN, LANG_ENGLISH},
	LANG_JAPANESE:   {LANG_JAPANESE, LANG_ENGLISH},
	LANG_BULGARIAN:  {LANG_BULGARIAN, LANG_ENGLISH},
	LANG_GEORGIAN:   {LANG_GEORGIAN, LANG_RUSSIAN, LANG_ENGLISH},
	LANG_NORWEGIAN:  {LANG_NORWEGIAN, LANG_ENGLISH},
	LANG_SWEDISH:    {LANG_SWEDISH, LANG_ENGLISH},
	LANG_CROATIAN:   {LANG_CROATIAN, LANG_ENGLISH},
	LANG_CHINESE:    {LANG_CHINESE, LANG_ENGLISH},
	LANG_PERSIAN:    {LANG_PERSIAN, LANG_ENGLISH},
	LANG_ROMANIAN:   {LANG_ROMANIAN, LANG_ENGLISH},
	LANG_HINDI:      {LANG_HINDI, LANG_ENGLISH},
	LANG_UKRAINIAN:  {LANG_UKRAINIAN, LANG_RUSSIAN, LANG_ENGLISH},
	LANG_MACEDONIAN: {LANG_MACEDONIAN, LANG_ENGLISH},
	LANG_SLOVENIAN:  {LANG_SLOVENIAN, LANG_ENGLISH},
	LANG_LATVIAN:    {LANG_LATVIAN, LANG_ENGLISH},
	LANG_SLOVAK:     {LANG_SLOVAK, LANG_ENGLISH},
	LANG_CZECH:      {LANG_CZECH, LANG_ENGLISH},
	LANG_AMHARIC:    {LANG_AMHARIC, LANG_ENGLISH},
}

var BLOGS_LANG = map[string]string{
	"laitman-ru":    "ru",
	"laitman-com":   "en",
	"laitman-es":    "es",
	"laitman-co-il": "he",
}
