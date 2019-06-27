namespace java com.github.jsix.go2o.rpc
namespace csharp com.github.jsix.go2o.rpc
namespace go go2o.core.service.auto_gen.rpc.foundation_service
include "ttype.thrift"


// 单点登录应用
struct SSsoApp{
    // 编号
    1: i32 ID
    // 应用名称
    2: string Name
    // API地址
    3: string ApiUrl
    // 密钥
    4: string Token
}

/** 行政区域 */
struct SArea  {
    1:i32 Code
    2:i32 Parent
    3:string Name
}


// 基础服务
service FoundationService{

   /** 获取注册表键值 */
   string GetRegistry(1:string key)
   /** 获取键值存储数据字典 */
   map<string,string> GetRegistries(1:list<string> keys)
   /** 创建自定义注册表项,@defaultValue 默认值,如需更改,使用UpdateRegistry方法  */
   ttype.Result CreateUserRegistry(1:string key,2:string defaultValue,3:string description)
   /** 更新注册表键值 */
   ttype.Result UpdateRegistry(1:map<string,string> registries)

   // 格式化资源地址并返回
   string ResourceUrl(1:string url)
   // 设置键值
   ttype.Result SetValue(1:string key,2:string value)
   // 删除值
   ttype.Result DeleteValue(1:string key)
   // 获取键值存储数据
   list<string> GetRegistryV1(1:list<string> keys)
   // 根据前缀获取值
   map<string,string> GetValuesByPrefix(1:string prefix)
   // 注册单点登录应用,返回值：
   //   -  1. 成功，并返回token
   //   - -1. 接口地址不正确
   //   - -2. 已经注册
   string RegisterApp(1:SSsoApp app)
   // 获取应用信息
   SSsoApp GetApp(1:string name)
   // 获取单点登录应用
   list<string> GetAllSsoApp()
   // 验证超级用户账号和密码
   bool SuperValidate(1:string user,2:string pwd)
   // 保存超级用户账号和密码
   void FlushSuperPwd(1:string user,2:string pwd)
   // 创建同步登录的地址
   string GetSyncLoginUrl(1:string returnUrl)
   // 获取地区名称
   list<string> GetAreaNames(1:list<i32> codes)
   // 获取下级区域
   list<SArea> GetChildAreas(1:i32 code)
}



