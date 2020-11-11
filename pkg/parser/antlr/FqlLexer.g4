lexer grammar FqlLexer;

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
OpenBracket: '[';
CloseBracket: ']';
OpenParen: '(';
CloseParen: ')';
OpenBrace: '{';
CloseBrace: '}';

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
MinusMinus: '--';
PlusPlus: '++';

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

// Group operators
Into: 'INTO';
Keep: 'KEEP';
With: 'WITH';
Count: 'COUNT';
All: 'ALL';
Any: 'ANY';
Aggregate: 'AGGREGATE';

// Unary operators
Like: 'LIKE';
Not: 'NOT' | '!';
In: 'IN';
Do: 'DO';
While: 'WHILE';

// Literals
Param: '@';
Identifier: Letter+ (Symbols (Identifier)*)* (Digit (Identifier)*)*;
StringLiteral: SQString | DQSring | BacktickString | TickString;
IntegerLiteral: [0-9]+;
FloatLiteral
    : DecimalIntegerLiteral Dot [0-9]+ ExponentPart?
    | DecimalIntegerLiteral ExponentPart?
    ;

NamespaceSegment: Identifier NamespaceSeparator;

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
fragment Symbols: '_';
fragment Digit
    : '0'..'9'
    ;
fragment DQSring: '"' ( '\\'. | '""' | ~('"'| '\\') )* '"';
fragment SQString: '\'' ('\\'. | '\'\'' | ~('\'' | '\\'))* '\'';
fragment BacktickString: '`' ('\\`' | ~'`')* '`';
fragment TickString: '´' ('\\´' | ~'´')* '´';
fragment NamespaceSeparator: '::';