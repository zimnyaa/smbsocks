#include <windows.h>
#include <stdio.h>

void dialpipe(char *);

BOOL WINAPI DllMain(
    HINSTANCE _hinstDLL,  
    DWORD _fdwReason,     
    LPVOID _lpReserved)   
{
    switch (_fdwReason) {
    case DLL_PROCESS_ATTACH:
		    
        break;
    case DLL_PROCESS_DETACH:
        
        break;
    case DLL_THREAD_DETACH:
        
        break;
    case DLL_THREAD_ATTACH:
		
        break;
    }
    return TRUE; 
}

typedef int (*goCallback)(const char*, int);
typedef struct _goArgs {
    char* argsBuffer;
    int bufferSize;
} goArgs;

void callGo(LPVOID param) {
    //goArgs args = *((goArgs*)param);
    //printf("ext entrypoint:\n\tstruct p:%llx\n\targs p:%llx\n\targs:%s\n\targs len:%x\n", args, args.argsBuffer, args.argsBuffer, args.bufferSize);
    dialpipe(param);
}


__declspec(dllexport) int __stdcall entrypoint(char *argsBuffer, int bufferSize, goCallback callback) {
    goArgs args;
    args.argsBuffer = argsBuffer;
    args.bufferSize = bufferSize;

    CreateThread(NULL, NULL, callGo, argsBuffer, NULL, NULL);

    callback("socks thread created.", 22);
    return 0;
}





