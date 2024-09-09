<?php

class Router
{

    public static $file = '../src/router/Routers.go';
    public static $sqlFile = 'menu.json';

    /**
     * 初始化一些路由
     */
    public static function init()
    {
        if (!is_file(self::$file)) {
            echo "文件定位失败: ", self::$file, "\n";
            exit;
        }

        $content = file_get_contents(self::$file);
        $lines = explode("\n", $content);
        $lineArr = [];
        foreach ($lines as $line) {
            $realLine = trim($line);
            if ($realLine == '' || strpos($realLine, '//') == 0 || (strpos($realLine, 'GET') === false && strpos($realLine, 'POST') === false)) {
                continue;
            }

            $comment = '缺少方法说明';
            $cArr = explode('//', $realLine);
            if (count($cArr) >= 2) {
                $comment = trim($cArr[1]);
            }

            $url = '';
            $sArr = explode(',', $cArr[0]);
            $pArr = explode('(', $sArr[0]);
            $url = trim($pArr[1], '"');

            $obj = (object) [
                'url'       => $url, //路由
                'comment'   => $comment, //备注
                'level'     => 1, //菜单级别
                'parent_id' => 0, //上级菜单
            ];

            $lineArr[] = $obj;

            print_r($lineArr);
        }
    }

    /**
     * 生成SQL
     */
    public static function sql()
    {
        $content = file_get_contents(self::$sqlFile);
        $menus = json_decode($content);

        $sqlFormat = "INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (%d, %s, '%s', '%s', '%d', '%s');";
        $base = 100000;
        echo "-- //", "\n";
        echo "-- * 严重警告: 此文件为只读, 请勿修改此文件", "\n";
        echo "-- * 如果要添加、修改菜单/结构, 请修改menu.json文件, 并将菜单结构反映其中", "\n";
        echo "-- * 修改完成 menu.json 文件之后, 执行 php -f ./create_menus.php > menus.sql 生成新的菜单文件即可", "\n";
        echo "-- //", "\n";
        echo "\n";
        echo "-- 清空当前菜单记录", "\n";
        echo "truncate menus;", "\n";
        echo "\n";
        foreach ($menus as $k1 => $menu1) {
            $menu1_id = $base * ($k1 + 1) + $k1;
            // 先写一级菜单
            echo "\n";
            $sql = sprintf($sqlFormat, $menu1_id, 0, $menu1->name, '#', $menu1->level, $menu1->icon);
            echo "-- 一级菜单: ", $menu1->name, "\n";
            echo $sql, "\n";
            foreach ($menu1->menus as $k2 => $menu2) {
                $menu2_id = $menu1_id + (10000 * ($k2 + 1)) + $k2;
                $sql = sprintf($sqlFormat, $menu2_id, $menu1_id, $menu2->name, $menu2->url, $menu2->level, '');
                echo "\t-- 二级菜单: ", $menu2->name, "\n";
                echo "\t", $sql, "\n";
                foreach ($menu2->menus as $k3 => $menu3) {
                    $menu3_id = $menu2_id + (1000 * ($k3 + 1)) + $k3;
                    $sql = sprintf($sqlFormat, $menu3_id, $menu2_id, $menu3->name, $menu3->url, $menu3->level, '');
                    echo "\t\t", $sql, "\n";
                    foreach ($menu3->menus as $k4 => $menu4) {
                        $menu4_id = $menu3_id + (100 * ($k4 + 1)) + $k4;
                        $sql = sprintf($sqlFormat, $menu4_id, $menu3_id, $menu4->name, $menu4->url, $menu4->level, '');
                        echo "\t\t\t", $sql, "\n";
                        if (!isset($menu4->menus) || count($menu4->menus) == 0) {
                            continue;
                        }
                        foreach ($menu4->menus as $k5 => $menu5) {
                            $menu5_id = $menu4_id + (10 * ($k5 + 1)) + $k5;
                            $sql = sprintf($sqlFormat, $menu5_id, $menu4_id, $menu5->name, $menu5->url, $menu5->level, '');
                            echo "\t\t\t\t", $sql, "\n";
                            foreach ($menu5->menus as $k6 => $menu6) {
                                $menu6_id = $menu5_id + ($k6 + 1);
                                $sql = sprintf($sqlFormat, $menu6_id, $menu5_id, $menu6->name, $menu6->url, $menu6->level, '');
                                echo "\t\t\t\t\t", $sql, "\n";
                            }
                        }
                    }
                }
            }
        }
    }

    /**
     * 递归处理所有sql
     * @param array $menus
     * @param array $urls
     */
    private static function getUrls($menus, &$urls)
    {
        foreach ($menus as $menu) {
            if (isset($menu->url)) {
                $urls[$menu->url] = 1;
            }
            if (isset($menu->menus) && count($menu->menus) > 0) {
                self::getUrls($menu->menus, $urls);
            }
        }
    }

    /**
     * 检测路由缺失状态
     */
    public static function check()
    {
        // 从json当中获取所有url
        $realRouters = [];
        $menus = json_decode(file_get_contents(self::$sqlFile));
        self::getUrls($menus, $realRouters);
        $realRouters = array_keys($realRouters);

        // 从Routers.go 读取所有真实的URL
        $lines = explode("\n", file_get_contents(self::$file)); // 拆分文件行
        $urls = [];
        foreach ($lines as $line) {
            $matched = NULL;
            $realLine = trim($line);
            if (preg_match('/"\/([a-z_]+[^"]*)"/', $realLine, $matched)) {
                $url = trim($matched[0], '"');
                $urls[$url] = 1;
            }
        }
        $urls = array_keys($urls);

        foreach ($urls as $url) {
            if (isset($realRouters[$url])) {
                echo $url, "\n";
            }
        }
    }

    /**
     * 运行
     * @param array $args
     */
    public static function run($args)
    {
        if (count($args) > 1 && trim($args[1]) == 'check') {
            self::check();
            return;
        }
        self::sql();
    }
}

Router::run($argv);
