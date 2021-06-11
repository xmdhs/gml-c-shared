//仅作注释用，不要导入此文件

typedef struct
{
    //游戏路径，需要创建一个 .minecraft 文件夹，例如 D:/mc/.minecraft
    char *Minecraftpath; 
    //分配内存大小，单位为 mb
    long long RAM;
    //玩家名
    char *Name;
    //uuid
    char *UUID;
    //正版/外置登录相关，离线登录可随便设置
    char *AccessToken;
    //要启动的游戏版本，可用 ListVersion 方法找到被识别的版本
    char *Version;
    //外置登录的 api 地址。
    char *ApiAddress;
    //自定义的 java 参数，字符串数组，可以通过 NewChar 和 SetChar 来构建
	char **Flag;
    //字符串数组的长度
	int  flag_len;
    //是否开启外置登录
	int  independent;
} Gameinfo;

typedef struct
{
    //code 对应的信息见 err.md
    int  code;
    //详细的错误信息
    char *msg;
} err;


typedef struct
{
    //字符串数组
    char **charlist;
    //字符串数组的长度
    int len;
    //错误
	err e;

} GmlReturn;

//当下载文件时，发生错误时将调用此函数。
typedef void (*Fail)(char*);

//下载文件成功后，会调用此函数
//第一个参数表示下载成功的类型
/*第一个参数见下方表格解释。第二个参数传递剩余没下载的文件数量，成功下载一次，调用一次。

code | msg
-----|-------
1    | 下载游戏核心
2    | 下载资源文件
3    | 下载游戏库
*/
typedef void (*Ok)(int,int);

typedef struct
{
    char *Username;
    char *ClientToken;
    char *UUID;
    char *AccessToken;
    char *ApiAddress;
} AuthDate;

typedef struct
{
    char *Username;
    char *UUID;
    char *AccessToken;
} MsAuthDate;

#ifdef __cplusplus
extern "C" {
#endif

//获取字符串数组中的字符串
extern char* Getchar(char** charlist, long long int index);
//释放字符串数组的内存，下面的函数返回的字符串数组都是分配在堆上的，需要用此函数释放
extern void Freechar(char** charlist, long long int len);
//创建一个字符串数组
extern char** NewChar(long long int len);
//设置字符串数组指定的位置的字符串
extern void SetChar(char** cc, long long int index, char* achar);
//生成启动游戏需要的参数
extern GmlReturn GenCmdArgs(Gameinfo g);
//设置下载文件和正版登录的代理，非线程安全，全局共用
extern err SetProxy(char* httpProxy);
//下载游戏
//version 下载版本，可通过 ListDownloadVersion 查找可下载的版本
//Type 下载使用的下载源，留空将按照权重的随机使用三个下载源，也可以自行设置。例如 vanilla|bmclapi 表示随机使用原版下载源和 bmclapi 下载源。mcbbs 表示只使用 mcbbs 下载源
//Minecraftpath 下载的路径，例如 D:/mc/.minecraft
//downInt 下载使用的协程数，通常 64 即可，因为每一个文件只使用了一个协程，多了虽然性能上没啥问题但是没意义。
extern err Download(char* version, char* Type, char* Minecraftpath, int downInt, Fail fail, Ok ok);
//检查游戏的完整性，第一次某个版本时，必须检查一次。
extern err Check(char* version, char* Type, char* Minecraftpath, int downInt, Fail fail, Ok ok);
//列出可启动版本
//path 例如 D:/mc/.minecraft/version
extern GmlReturn ListVersion(char* path);
//查看可下载游戏版本类型，比如正式版之类的。 type 参数和 Download 里的意义一样
extern GmlReturn ListDownloadType(char* Type);
//查看此类型中所有的可下载的版本
extern GmlReturn ListDownloadVersion(char* VerType, char* Type);

/* Return type for Auth */
struct Auth_return {
	AuthDate r0;
	err r1;
};
//外置登录和正版登录
//clientToken 客户端 id，可随机生成，需保证每个用户对应的 ClientToken 是不变的，否则会要求重新登录。建议直接 md5 用户名就行。
//ApiAddress 可不输入完整的 api 地址，会按照协议补全，如果正版登录，则无需设置此项。
extern struct Auth_return Auth(char* ApiAddress, char* username, char* email, char* password, char* clientToken);
//验证 AccessToken 的有效性，建议每次启动游戏前，都验证一次
extern err Validate(char* AccessToken, char* ClientToken, char* ApiAddress);

/* Return type for Refresh */
struct Refresh_return {
	AuthDate r0;
	err r1;
};
//若 Validate 验证 AccessToken 失效，则可通过此方法刷新，如果刷新无效，则重新登录就是。
extern struct Refresh_return Refresh(char* AccessToken, char* ClientToken, char* ApiAddress);

/* Return type for MsAuth */
struct MsAuth_return {
	MsAuthDate r0;
	err r1;
};
//微软登录，将弹出一个窗口让玩家在其中进行登录，需要安装 chrome，或者新版 edge（windows only）
extern struct MsAuth_return MsAuth();

/* Return type for MsAuthValidate */
struct MsAuthValidate_return {
	MsAuthDate r0;
	err r1;
};
//验证微软登录的 AccessToken 有效性，建议首次登录后也使用此方法验证，因为即使没有购买游戏，MsAuth 也会返回有效的内容。
extern struct MsAuthValidate_return MsAuthValidate(char* AccessToken);

#ifdef __cplusplus
}
#endif
