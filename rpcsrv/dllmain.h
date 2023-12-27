#include <windows.h>

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

