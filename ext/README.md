#	Markdown


##	Link


http://mdblog.daoapp.io/

### [Some **bold** _italic_ link](http://mdblog.daoapp.io/)



##	Font

This is a _test_

This is a *test*

This is a **test**

~~cross out~~

>Golang is good.



## Table

Name        | Age
------------|------
Bob     	| 27
Alice   	| 23



##	Code

``` go
func getTrue() bool {
    return true
}
```


##	List

- [ ] This is an incomplete task.
- [x] This is done.

*   Red
*   Green
*   Blue

等同于：

+   Red
+   Green
+   Blue

也等同于：

-   Red
-   Green
-   Blue

有序列表则使用数字接着一个英文句点：

1.  Bird
2.  McHale
3.  Parish


*   Lorem ipsum dolor sit amet, consectetuer adipiscing elit.
Aliquam hendrerit mi posuere lectus. Vestibulum enim wisi,
viverra nec, fringilla in, laoreet vitae, risus.
*   Donec sit amet nisl. Aliquam semper ipsum sit amet velit.
Suspendisse id sem consectetuer libero luctus adipiscing.
*   A list item with a blockquote:

    > This is a blockquote
    > inside a list item.
*   一列表项包含一个列表区块：

		<代码写在这>

* * *

***

*****

- - -

---------------------------------------

This is [an example][id] reference-style link.
你也可以选择性地在两个方括号中间加上一个空格：

This is [an example] [id] reference-style link.
接着，在文件的任意处，你可以把这个标记的链接内容定义出来：

[id]: http://example.com/  "Optional Title Here"


I get 10 times more traffic from [Google] [1] than from
[Yahoo] [2] or [MSN] [3].

  [1]: http://google.com/        "Google"
  [2]: http://search.yahoo.com/  "Yahoo Search"
  [3]: http://search.msn.com/    "MSN Search"

![Alt text](/path/to/img.jpg)

![Alt text](/path/to/img.jpg "Optional title")