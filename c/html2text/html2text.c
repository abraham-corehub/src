#include "html2text.h"

/*  who is who?
 *  <b>Hello</b>
 *  <	: tag start
 *  b 	: tag name
 *  >	: tag end
 *  <b>	: Start Tag
 *  </b>: End Tag
 *  Hello: Text
 */


// General Functions
//--------------------------------------------------------------------------------------------//

static void die(int code,const char *msg, const char *arg1)
{	
	fprintf(stdout,"%s: %s %s\n", appname, msg,arg1?arg1:"");
	exit(code);
}

static void show_usage(char *more){
	fprintf(stdout,"%s %s - Simple HTML to Text converter\nBuild : %s %s\n\n",appname,appver,__DATE__,__TIME__);
	if(more)
		fprintf(stdout,"\n%s\n\n",more);
	fprintf(stdout,"Usage: %s html_file [ -o output_file.txt ]\n",appname);
	exit(1);
}

static void putchars(char *p,int n){
	char *e=p+n;
	while (*p && p<e){
		if(*p>' ')
			fputc(*p,stderr);
		else
			fputc('?',stderr);
		p++;
	}
	fputs(".\n",stderr);
}

// Parsing Helper Functions
//--------------------------------------------------------------------------------------------//

/* Just using std c functions */
#define is_alphanum isalnum
#define is_digit isdigit
#define is_alpha isalpha

/* Checks if charecter is a white space */
static inline int is_space(char c){
	return (c==' ' || c=='\n' || c=='\r' || c=='\t' || c=='\v' || c=='\f');
}

/* Returns the tagname length */
static int tagname_len(char *p){
	char *c=p;	
	if (is_alpha(*c))
		while(is_alphanum(*c))
			c++;		
	return c-p;
}

/* Checks if the given endtag is valid */
static inline int is_comment(char *p){
	return *p=='!' || *p=='?' || *p=='%';
}

/* Checks if the given tagname is valid */
static int is_tagname(char *p){
	if(*p=='/') p++;
	return is_alpha(*p) || is_comment(p);	
}

/* Checks if the given tag is valid */
static int is_tag(char *p){
	//return *p=='<' && !is_space(*++p);	
	if(*p=='<')
		return is_tagname(++p);
	return 0;
}

/* Checks if the given endtag is valid */
static inline int is_endtag(char *p){
	return *p=='<' && *++p=='/' && is_tagname(++p);
}

/* is_block_tag ? */
static int is_block_tag(char *p){	
	int i;	
	for(i=0;i<BLOCK_TAG_COUNT;i++){						
		if(*BLOCK_TAG[i]=='/' || *BLOCK_TAG[i]=='<'){ // forced start/end tags in list
			if(!strcasecmp(p-1,BLOCK_TAG[i]))	return -1;
		}else{
			if(!strcasecmp(p,BLOCK_TAG[i]))	return -1;
		}
	}
	return 0;
}

/* is_special_tag (style, script...) ? */
static int is_special_tag(char *p){	
	int i;
	for(i=0;i<SPECIAL_TAG_COUNT;i++){		
		if(!strcasecmp(p,SPECIAL_TAG[i]))
			return 1;		
	}
	return 0;
}

/* Returns pointer to char after needle from haystack */
static char* skip_past(char* needle,char *haystack, char *end){
	//printf("Skipping : %s\n",needle);
	char *p=haystack, *np=needle;
	int nlen=strlen(needle);
	while(*p && p<end){
 		if(tolower(*p)==tolower(*np))
 			np++;
 		else if(np!=needle){	 						
 			p=p-(np-needle);	 				
 			np=needle;	 
 		}
 		p++;
 		if(!*np)
 			return p;
	}
	return end;	
}

/* Skips to end of tag, ie > */
static char * skip_tag(char *p,char *e){	
	while(p++<e){
		if(*p=='"')
			p=skip_past("\"",p+1,e);
		else if	(*p=='\'')
			p=skip_past("\'",p+1,e);
		if(*p=='>')
			return p++;
	}
	return e;	
}

/* Skips Comments/Serverside tags */
static char * skip_comment(char *p,char *e){	
	char c;
	if(*p=='!'){			// <!> or <!-- -->
		if(*++p=='-' && *(p+1)=='-')
			return skip_past("-->",p,e)-1;
		else
			return skip_past(">",p,e);
	}else{					// <? ?> or <% %> 
		c=*p;
		while(*(p-1)!=c)
			p=skip_tag(p,e);
		return p;
	}	
}


/* Writes a charecter to output */
static void write_char(char c, ParserState *pstate){
	//if(pstate->in_body || pstate->in_title)
		fputc(c,pstate->file);
	pstate->newline=(c=='\n');
	if(pstate->in_title)
		pstate->title_len+=1;
}

/* Writes a string to output */
static void write_str(const char *str, ParserState *pstate){
	for(;*str;str++)
		write_char(*str,pstate);
}

/* takes care of actions needed for tags */
static char* deal_tags(char *p,char *end, ParserState *pstate){
	char tagsep,*te;
	
	int len=tagname_len(p);		
	if(len==0) return p; 	//invalid tagname, treat as text	
	
	te=p+len; tagsep=*te; 	//bakup
	*te='\0'; 		// insert null after tag to make it a C string 	
	//printf("\n------<%s>------[%d]\n",p,len);
	
	if(0);
	else if(!strcasecmp(p,	"br" 	))	do_tag_br(p,pstate);
	else if(!strcasecmp(p,	"hr" 	))	do_tag_hr(p,pstate);	
	else if(!strcasecmp(p,	"body"	))	do_tag_body(p,pstate);
	else if(!strcasecmp(p,	"title" ))	do_tag_title(p,pstate);			
	else if(!strcasecmp(p,	"pre" 	))	p=do_tag_pre(p,pstate);
	else if(is_block_tag(p))	do_block_tag(p,pstate);
	else if(is_special_tag(p)!=0){ //<script>, <style>...
		p--;
		*p='/';		
		p=skip_past(p,te+1,end)-2;
	}		
	*te=tagsep;
	
	return skip_tag(p,end);
}

/* takes care of actions needed for end tags */
static char* deal_etags(char *p,char *end, ParserState *pstate){
	char tagsep,*te;
	
	int len=tagname_len(p);		
	if(len==0) return p;	//invalid tagname, treat as text
	te=p+len;
	tagsep=*te; 	//bakup
	*te='\0'; 		// insert null after tag to make it a C string 	
	//printf("\n------</%s>------[%
	if(0);	
	else if(!strcasecmp(p,	"body"	))	do_etag_body(p,pstate);
	else if(!strcasecmp(p,	"title"	))	do_etag_title(p,pstate);
	else if(!strcasecmp(p,	"pre" 	))	do_etag_pre(p,pstate);
	else if(is_block_tag(p))			do_block_etag(p,pstate);

	*te=tagsep;
	return skip_tag(p,end);
}

/* translates &amp; to & and so */
static char* deal_entities(char *p, ParserState *pstate){	
	int i,elen=strchr(p,';')-p;
	char t;
	if(elen > 2 && elen<10){
		t=p[elen+1];		
		p[elen+1]=0;
		for(i=0;i<ENTITY_COUNT;i+=3){
			if((!strcasecmp(p,ENTITY[i+1])) || !strcasecmp(p,ENTITY[i+2])){
				write_str(ENTITY[i],pstate);
				pstate->hitspace=0;
				p[elen+1]=t;
				return p+elen;	
			}
		}
		p[elen+1]=t;
	}
	return 0;
}

/* does the parsing job */
static void html2txt(char *html,size_t length, FILE *outfile){
	char *p=html, *end=(html+length),*t;
	ParserState ps={0};
	ps.hitspace=1;
	ps.file=outfile;
	
	while(p<end){
		if(is_tag(p)){
			p++;
			if(is_comment(p))									//Comments
				p=skip_comment(p,end);				
			else 
			if(*p=='/'){										//End Tags
				p++;
				p=deal_etags(p,end,&ps);
			}else												//Tags
				p=deal_tags(p,end,&ps);			
		}else{						
	 		if(is_space(*p) && !ps.in_pre){						//Space
	 			if(!ps.hitspace){
	 				write_char(' ',&ps);
	 				ps.hitspace=1;
	 			}		
	 		}else{
	 		 	if((*p=='&') && (t=deal_entities(p,&ps)))		//HTML Entity					
	 		 		p=t;
	 		 	else
	 		 		write_char(*p,&ps);							//Text	 
	 		 	ps.hitspace=0;
	 		}
		}
		if(ps.done)
			break;
		p++;
	}
}

// Do actions for specific tags
//--------------------------------------------------------------------------------------------//

void static inline do_putc(char c,ParserState *pstate){
	fputc(c,pstate->file);	
}

void static inline do_newline(ParserState *pstate){
	if(!pstate->newline)
		write_char('\n',pstate);
}

static void do_block_tag(char *p, ParserState *pstate){
	do_newline(pstate);
	pstate->hitspace=1;
}

static void do_block_etag(char *p, ParserState *pstate){
	do_newline(pstate);
	pstate->hitspace=1;
}

static void do_tag_br(char *p, ParserState *pstate){
	do_newline(pstate);
	pstate->newline=0;
}

static void do_tag_hr(char *p, ParserState *pstate){
	int i;
	do_newline(pstate);		
	for(i=79;i;i--)
		do_putc('_',pstate);
	do_putc('\n',pstate);
	pstate->hitspace=1;	
}

static void do_tag_title(char *p, ParserState *pstate){
	if(pstate->in_body) return;
	pstate->in_title=-1;	
	pstate->title_len=0;
}
static void do_etag_title(char *p, ParserState *pstate){
	if(pstate->in_body) return;
	pstate->in_title=0;
	do_putc('\n',pstate);
	if(is_space(*(p-3))) pstate->title_len--; //trailing space
	if(pstate->title_len>79) pstate->title_len=80; //length 80 cap
	for(;pstate->title_len;pstate->title_len--)
		do_putc('=',pstate);
	write_str("\n\n\n",pstate);
	pstate->hitspace=1;
}

static void do_tag_body(char *p, ParserState *pstate){
	pstate->in_body=1;
	pstate->hitspace=1;	
}
static void do_etag_body(char *p, ParserState *pstate){
	pstate->in_body=0;
	pstate->done=1;	
	do_putc('\n',pstate);
}
static char* do_tag_pre(char *p, ParserState *pstate){
	pstate->in_pre=1;
	pstate->hitspace=0;
	while(is_space(*p))
		p++;
	return p;
}
static void do_etag_pre(char *p, ParserState *pstate){
	pstate->in_pre=0;	
	pstate->hitspace=1;
	do_putc('\n',pstate);
}


/* Read the whole file to mem and feed it to html2txt() */
int main(int argc, const char *argv[])
{
	FILE *infile,*outfile;
	size_t infile_size,read;
	char *buf;
	if(argc<2)
		show_usage(0);
	
	infile=fopen(argv[1],"r");
	if(!infile)
		die(0,"Failed to Open Input HTML File",argv[1]);
	fseek(infile, 0, SEEK_END);
	infile_size=ftell(infile);
	rewind(infile);
	buf=malloc(infile_size * sizeof(char)+1);
	if(!buf)
		die(2,"Failed to allocate infile buffer.",0);
	read=fread(buf, sizeof(char), infile_size, infile);
	fclose(infile);
	buf[read]=0;
	
	if(argc==4 && argv[2][0]=='-' && argv[2][1]=='o'){
		outfile=fopen(argv[3],"w");
		if(!outfile)
			die(0,"Failed to open output file",argv[2]);
	}else
		outfile=stdout;

	
	html2txt(buf,read,outfile);		
	

	free(buf);
	if(outfile!=stdout) fclose(outfile);
	return 0;
}
