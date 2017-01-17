#	[mdbgEg](https://github.com/toukii/mdbgEg)

**markdown blog engine for github**

***

>	USAGE

>>	go get github.com/toukii/mdbgEg

>>	$GOPATH/bin/mdbgEg

>	visit:  _http://localhost/_

##	BUILD

 说明：该项目是一款静态博客引擎，从支持git的网站如[GITHUB](https://github.com/)上pull内容，将markdown文件转为html文件，并提供该html的文件服务，项目参考的技术详见[参考](http://mdblog.daoapp.io/Item/mdbg/) 。

以下是运行该项目需要的文件：

*	Dockerfile，基于[DaoCloud](https://daocloud.io)的CI，方便部署。Dockerfile格式如下(**只需将user/repo替换为自己的项目即可。**)：

>		FROM golang

>		WORKDIR /app/gopath/mdblog
>		ENV GOPATH /app/gopath

>		RUN git clone --depth 1 git://github.com/**user/repo**.git .

>		RUN go get -u github.com/toukii/mdbgEg
>		RUN mv $GOPATH/bin/mdbgEg /bin/mdbgEg

>		EXPOSE 80

>		CMD ["/bin/mdbgEg"]


*	添加WebHook，在自己的GITHUB（或其他支持webhook）项目中添加回调接口，如网站地址为 **http://blog.daoapp.io**，添加的接口为 **http://blog.daoapp.io/callback**  ;webhook的类型选取push即可。回调函数请求中包含了修改了哪些文件，如图所示：

	![callback](http://7xku3c.com1.z0.glb.clouddn.com/callback.png)


*	可以在自己的项目根目录添加文件名为**theme.thm**的模板文件:

		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
				<title>mdblog</title>
			</head>
			<body>
				<div>
					{{.MDContent}}
				</div>
			</body>
		</html>

**完成以上步骤，提交"代码"即可完成更新;另外，也支持在线更新模板。**

*	![kiss me](http://7xku3c.com1.z0.glb.clouddn.com/baby.gif)

_2015-11-29_