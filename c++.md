教材[[C++ Primer (第5版中文版) .pdf]]

# 补充知识

## argc和argv参数？
~~~c++
int main (int argc,char *argv[])
{
	return 0;
	
}
~~~

==argc==:
	The first parameter, argc (argument count) is an integer that indicates how many arguments were entered on the command line when the program was started
==argv==:
	an array of pointers to arrays of character objects. The array objects are null-terminated strings, representing the arguments that were entered on the command line when the program was started.

