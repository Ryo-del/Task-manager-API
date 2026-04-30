loHttp = CREATEOBJECT("MSXML2.ServerXMLHTTP.6.0")

lcUrl = "https://jsonplaceholder.typicode.com/todos/1" &&ссылка на backend


loHttp.Open("Get", lcUrl, .F.)

loHttp.Send()

lcResponse = loHttp.ResponseText

_SCREEN.Cls
ACTIVATE SCREEN

? "STATUS: ", loHttp.Status
? "RESPONSE BODY:"
? lcResponse

WAIT WINDOW "Готово!"