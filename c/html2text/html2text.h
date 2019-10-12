/*
 * html2txt - Simple HTML to Text converter. 
 *
 * @author: Nahar < http://scr.im/nahar >
 *
 */
 
 static const char
    appname[]   = "html2txt",
    appver[]	= "1.0";

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <ctype.h> //isalpha, isspace....
#include <string.h>

//Fix MSVC
#ifdef _MSC_VER
  #define inline __forceinline
  #define strcasecmp _stricmp
#endif


typedef struct _ParserState {
	// Flags
	signed in_title: 1;
    signed in_body : 1;
    signed in_pre  : 1;
    signed hitspace: 1;
    signed newline : 1;
	signed done	   : 1;

    int title_len;
    FILE * file;    
} ParserState;

const char BLOCK_TAG[][15]={ 
	#include "block_tags.txt" 
};
const int  BLOCK_TAG_COUNT = sizeof(BLOCK_TAG)/15;

const char SPECIAL_TAG[][10]={
	#include "special_tags.txt" 
};
const int  SPECIAL_TAG_COUNT = sizeof(SPECIAL_TAG)/10;

const char ENTITY[][10]={ 
	#include "entities.txt" 
};
const int ENTITY_COUNT= sizeof(ENTITY)/10;


static void do_tag_hr(char *p, ParserState *pstate);
static void do_tag_br(char *p, ParserState *pstate);
static void do_tag_title(char *p, ParserState *pstate);
static void do_etag_title(char *p, ParserState *pstate);
static void do_tag_body(char *p, ParserState *pstate);
static void do_etag_body(char *p, ParserState *pstate);
static char* do_tag_pre(char *p, ParserState *pstate);
static void do_etag_pre(char *p, ParserState *pstate);

static void do_block_tag(char *p, ParserState *pstate);
static void do_block_etag(char *p, ParserState *pstate);
