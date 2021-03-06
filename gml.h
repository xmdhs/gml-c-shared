//仅作注释用，不要导入此文件


//char* 均使用 utf-8 编码

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
    //code 对应的信息见下方注释
    int  code;
    //详细的错误信息
    char *msg;
} err;

/*code | msg
-----|----------------------------------------------
-1   | 未知错误
1    | 文件不存在
2    | json 错误
3    | minecraft json 格式错误
4    | 找不到该版本
5    | 下载失败次数过多
6    | authlib 登录时，有多个账户，且并没有选择（也就是没有设置 username 参数）
7    | 没有这个角色
8    | 通常是密码错误，或者登录过期
9    | authlib 找不到可用档案，也就是没有创建角色之类的
10   | accessToken 失效
11   | 登录微软账户时打开的浏览器中，打开了其他页面
12   | 尝试重新登录微软账户
13   | 没有购买游戏或者没有迁移账号
14   | 没有安装 chrome 或者新版本的 edge（windows only）
*/

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

//下载或检查游戏，成功或失败时，将调用此函数
typedef void (*gmlfinish)(err);

typedef struct
{
    char *Username;
    char *ClientToken;
    char *UUID;
    char *AccessToken;
    char *ApiAddress;
    //可用的用户名，当 Auth 返回 6 错误时，可以选择此列表中的用户名，设置在 Auth 的 username 参数上，来选择档案。
    //只有外置登录有单账号，多档案。
    char **availableProfiles;
    //长度
    int availableProfilesLen;
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
//其实就是导出了 c 语言的 malloc 函数，为了避免变量名冲突，所以首字母大写了
extern void* Malloc(int i);
//生成启动游戏需要的参数
extern GmlReturn GenCmdArgs(Gameinfo g);
//设置下载文件和正版登录的代理，非线程安全，全局共用
extern err SetProxy(char* httpProxy);
//下载游戏
//version 下载版本，可通过 ListDownloadVersion 查找可下载的版本
//Type 下载使用的下载源，留空将按照权重的随机使用三个下载源，也可以自行设置。例如 vanilla|bmclapi 表示随机使用原版下载源和 bmclapi 下载源。mcbbs 表示只使用 mcbbs 下载源
//downInt 下载协程数，通常 64 即可，因为一个文件还是使用一个协程下载，多了没意义
//Minecraftpath 下载的路径，例如 D:/mc/.minecraft
//调用后将立刻返回一个 int64，可以使用 Cancel 函数，将此 int64 传入，取消下载操作
extern long long int Download(char* version, char* Type, char* Minecraftpath, int downInt, Fail fail, Ok ok, gmlfinish finish);
//检查游戏的完整性，第一次某个版本时，必须检查一次。
extern long long int Check(char* version, char* Type, char* Minecraftpath, int downInt, Fail fail, Ok ok, gmlfinish finish);
//取消下载或者检查游戏完整性
extern void Cancel(long long int id);
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
//username 外置登录时需要，具体见 AuthDate 处的注释
//不要保存密码，保存 AccessToken
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
