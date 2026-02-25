lexer grammar FqlLexer;

tokens { TemplateExprEnd }

// Skip
MultiLineComment: '/*' .*? '*/'             -> channel(HIDDEN);
SingleLineComment: '//' ~[\r\n\u2028\u2029]* -> channel(HIDDEN);
WhiteSpaces: [\t\u000B\u000C\u0020\u00A0]+ -> channel(HIDDEN);
LineTerminator: [\r\n\u2028\u2029] -> channel(HIDDEN);

// Punctuation
Colon: ':';
SemiColon: ';';
Dot: '.';
Comma: ',';
Tilde: '~';
OpenBracket: '[';
CloseBracket: ']';
OpenParen: '(';
CloseParen: ')';
OpenBrace: '{' { if len(l.templateDepth) > 0 { l.templateDepth[len(l.templateDepth)-1]++ } };
CloseBrace: '}' {
	if len(l.templateDepth) > 0 {
		l.templateDepth[len(l.templateDepth)-1]--
		if l.templateDepth[len(l.templateDepth)-1] == 0 {
			l.templateDepth = l.templateDepth[:len(l.templateDepth)-1]
			l.PopMode()
			l.SetType(FqlLexerTemplateExprEnd)
		}
	}
};

// Comparison operators
Gt: '>';
Lt: '<';
Eq: '==';
Gte: '>=';
Lte: '<=';
Neq: '!=';

// Arithmetic operators
Multi: '*';
Div: '/';
Mod: '%';
Plus: '+';
Minus: '-';
Increment: '++';
Decrement: '--';

// Logical operators
And: 'AND' | '&&';
Or: 'OR' | '||';

// Other operators
Range: Dot Dot;
Assign: '=';
QuestionMark: '?';
RegexNotMatch: '!~';
RegexMatch: '=~';

// Keywords
// Common Keywords
For: 'FOR';
Return: 'RETURN';
Waitfor: 'WAITFOR';
Dispatch: 'DISPATCH';
Options: 'OPTIONS';
Timeout: 'TIMEOUT';
Every: 'EVERY';
Backoff: 'BACKOFF';
Jitter: 'JITTER';
Exists: 'EXISTS';
Value: 'VALUE';
Throw: 'THROW';
Distinct: 'DISTINCT';
Filter: 'FILTER';
Sort: 'SORT';
Limit: 'LIMIT';
Let: 'LET';
Collect: 'COLLECT';
SortDirection: 'ASC' | 'DESC';
None: 'NONE';
Null: 'NULL';
BooleanLiteral: 'TRUE' | 'true' | 'FALSE' | 'false';
Use: 'USE';
As: 'AS';
At: 'AT';
Least: 'LEAST';

// Group operators
Into: 'INTO';
Keep: 'KEEP';
With: 'WITH';
All: 'ALL';
Any: 'ANY';
Aggregate: 'AGGREGATE';

// Wait operators
Event: 'EVENT';

// Unary operators
Like: 'LIKE';
Not: 'NOT' | '!';
In: 'IN';
Do: 'DO';
While: 'WHILE';
Step: 'STEP';

// Literals
Param: '@';
Identifier: Letter+ (Symbols (Identifier)*)* (Digit (Identifier)*)*;
IgnoreIdentifier: Underscore;
StringLiteral: SQString | DQSring | TickString;
BacktickOpen: '`' -> pushMode(TEMPLATE);
DurationLiteral
    : DecimalIntegerLiteral (Dot [0-9]+)? ExponentPart? DurationUnit
    ;
IntegerLiteral: [0-9]+;
FloatLiteral
    : DecimalIntegerLiteral Dot [0-9]+ ExponentPart?
    | DecimalIntegerLiteral ExponentPart?
    ;

NamespaceSegment: Identifier NamespaceSeparator;

UnknownIdentifier: .;

// Fragments
fragment HexDigit
    : [0-9a-fA-F]
    ;
fragment DecimalIntegerLiteral
    : '0'
    | [1-9] [0-9]*
    ;
fragment ExponentPart
    : [eE] [+-]? [0-9]+
    ;
fragment Letter
    : 'A'..'Z' | 'a'..'z'
    ;
fragment Symbols: Underscore;
fragment Underscore: '_';
fragment Digit
    : '0'..'9'
    ;
fragment DQSring: '"' ( '\\'. | '""' | ~('"'| '\\') )* '"';
fragment SQString: '\'' ('\\'. | '\'\'' | ~('\'' | '\\'))* '\'';
fragment TickString: '´' ('\\´' | ~'´')* '´';
fragment NamespaceSeparator: '::';
fragment DurationUnit
    : 'MS'
    | 'S'
    | 'M'
    | 'H'
    | 'D'
    ;

mode TEMPLATE;

TemplateExprStart
    : '${' {
		l.templateDepth = append(l.templateDepth, 1)
		l.PushMode(0)
	}
    ;

TemplateChars
    : ( '\\' . | ~[`\\$] | '$' { p.GetInputStream().LA(1) != '{' }? )+
    ;

BacktickClose
    : '`' -> popMode
    ;
