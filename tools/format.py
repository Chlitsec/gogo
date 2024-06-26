# -*- coding: utf-8 -*-
import json
import click
import zlib
from pathlib import Path
from ipaddress import ip_address
from itertools import groupby
from functools import partial


class ResultJson:
    def __init__(self, json: dict):
        self.json = json

    def __getattr__(self, item):
        # assert item in self.namemap.keys(), "no found key in " + " ".join(self.namemap.keys())
        return self.json.get(item, "")

    def __getitem__(self, item):
        return getattr(self, item)

    # @staticmethod
    # def trans(item):
    #     return ResultJson.namemap[item]


class GogoResult:
    def __init__(self, result: dict):
        self.result = ResultJson(result)

    def __eq__(self, other):
        return self.target == other.target

    def __hash__(self):
        return hash(self.target)

    def __getattr__(self, item):
        return self.result[item]

    def __getitem__(self, item):
        return getattr(self, item)

    @property
    def vulns(self):
        vs = []
        for k,v in self.result["vulns"].items():
            vs.append(f"[{k}: payload:{v.get('payload','')} detail:{v.get('detail','')}]")
        return " ".join(vs)

    @property
    def frameworks(self):
        frames = []
        for k, v in self.result["frameworks"].items():
            frames.append(f"{k}:{v.get('version','')}")
        return "|".join(frames)

    @property
    def app(self):
        return self.frameworks

    def first_framework_name(self):
        frameworks = self.result["frameworks"]
        if len(frameworks):
            return frameworks[0][ResultJson.trans("framework_name")]
        else:
            return ""

    @property
    def url(self):
        return f'{self.protocol}://{self.ip}:{self.port}'

    @property
    def target(self):
        return f'{self.ip}:{self.port}'

    def dict2string(self, l):
        if l:
            return "|".join(["".join(map(str, i.values())) for i in l]).lower()
        else:
            return ""

    def equal(self, key, value):
        if self[key] == value:
            return True
        else:
            return False

    def contain(self, key, value):
        if value in self[key]:
            return True
        else:
            return False

    def to_json(self):
        return json.dumps(self.result.json)

    def to_dict(self):
        return self.result.json

    def gets(self, *args):
        return [self[i] for i in args]

    def has_framework_name(self, name):
        for framework in self["frameworks"]:
            if name == framework["name"]:
                return True
        return False


class GogoResults:
    operatormap = {
        "!=": (True, True),
        "==": (True, False),
        "!:": (False, True),
        "::": (False, False)
    }

    brutable = {
        "mariadb": "MYSQL",
        "mysql": "MYSQL",
        "microsoft rdp": "RDP",
        # "oracle database": "ORACLE",
        "microsoft sqlserver": "MSSQL",
        "mssql": "MSSQL",
        "smb": "SMB",
        "redis": "REDIS",
        "vnc": "VNC",
        "elasticsearch": "ELASTICSEARCH",
        "postgreSQL": "POSTGRESQL",
        "mongo": "MONGO",
        "ssh": "SSH",
        "ftp": "FTP"
    }

    def __init__(self, results: list):
        self.config = results[0]
        self.results = [GogoResult(result) for result in results[1:len(results)-2]]

    def __getattr__(self, item):
        return [result[item] for result in self.results if result[item]]

    def __getitem__(self, item):
        return getattr(self, item)

    def __or__(self, other):
        self.results = list(set(self.results + other.results))
        return self

    def __iter__(self):
        for result in self.results:
            yield result

    def equal_results(self, key, value, isnot=False):
        return GogoResults([result.to_dict() for result in self.results if isnot ^ result.equal(key, value)])

    def contain_results(self, key, value, isnot=False):
        return GogoResults([result.to_dict() for result in self.results if isnot ^ result.contain(key, value)])

    def split_expr(self, expr):
        for key in ResultJson.json:
            if expr.find(key) == 0:
                off = len(key)
                operator = expr[off:off + 2]
                assert operator in self.operatormap, "Confirm operator format, :: or !: or != or =="
                return key, expr[off + 2:].strip('"'), *self.operatormap[operator]

        # raise AttributeError(f"key not in {', '.join(ResultJson.namemap.keys())}")

    def expr(self, expr):
        k, v, isequal, isnot = self.split_expr(expr)
        if isequal:
            return self.equal_results(k, v, isnot)
        else:
            return self.contain_results(k, v, isnot)

    def exprs(self, exprs, isor=False):
        if len(exprs) == 0:  # 没有过滤条件则返回自身
            return self

        if not isor:
            results = self
            for expr in exprs:
                results = results.expr(expr)
        else:
            results = GogoResults([])
            for expr in exprs:
                results = results | self.expr(expr)
        return results

    @property
    def json(self):
        return ""

    @property
    def raw_json(self):
        return json.dumps([i.to_dict() for i in self])

    @property
    def brute(self):
        return self.exprs([f"frameworks::{k}" for k in self.brutable.keys()], True)

    @property
    def zombie(self):
        return json.dumps([{"IP": result.ip,
                            "Port": int(result.port),
                            "Server": self.brutable[result.first_framework_name().lower()]
                            }
                           for result in self.brute])

    @property
    def full(self):
        sorted_results = sorted(self.results, key=lambda x: int(ip_address(x.ip)))
        group_results = groupby(sorted_results, lambda x: x.ip)
        s = ""
        for ip, results in group_results:
            s += "[+] %s\n" % ip
            for result in results:
                s += f"\t{result.url}\t{result.midware}\t{result.language}\t{result.frameworks}\t{result.host}\t[{result.status}] {result.title}\t{result.vulns}\n"
        return s

    def gets(self, *args):
        return [result.gets(*args) for result in self]

    def getstrs(self, *args):
        return ["\t".join(i) for i in self.gets(*args)]


def toFile(filename, content: str):
    with open(filename, "a+", encoding="utf-8") as f:
        f.write(content)


def fixjson(content:str):
    if content.startswith("{") and not content.endswith("]}"):
        print("auto fix json!!!")
        return content + "]}"
    return content


def decompress(bs):
    flatedict = bytes(', ":'.encode())
    com = zlib.decompressobj(-15,zdict=flatedict)
    return com.decompress(bs).decode()


def loadResult(file):
    filename = Path(file.name)
    content = file.read()
    try:
        if ".dat" in filename.suffix:
            content = decompress(content)
            # content = fixjson(content)
        else:
            content = content.decode()
            # content = fixjson(content)
        return [json.loads(i) for i in content.split("\n") if i]
    except Exception as e:
        print(str(e))
        return []


@click.command()
@click.argument("files", nargs=-1, type=click.File("rb"))
@click.option('--output', '-o', default="full", help='Output format.')
@click.option('--expr', '-e', "exprs", multiple=True, help='filter rules')
@click.option('--outfile', '-f', help='output file')
@click.option('--or', '-or', "_or", default=False, is_flag=True)
def main(files, output, exprs, outfile, _or):
    """    使用帮助:

    \b
    过滤规则-e e.g: !!!请使用小写字母!!!
    全等匹配: port==443
    模糊匹配: title::系统
    排除: protocol!=tcp
    模糊排除: protocol!:tcp
    允许使用多个filter器,例如 -e port==443 -e title::系统

    \b
    开启增量匹配 -or
    flag字段,存在-or参数则匹配方式从递减修改为增量

    \b
    输出规则-o   e.g:
    选择任意一字段输出: ip port uri os host title midware http_stat languate frameworks vulns protocol vuln_name vuln_payload vuln_detail framework_name framework_version url target
    特殊输出:
    json
    zombie
    \b
    输出到文件-f: e.g: -f res.txt, 如果不指定-o则输出到命令行

    \b
    example:
    # 输出端口为443的ip
    python format.py input.json -e port==443
    \b
    # 输出端口为443的target结果
    python format.py input.json -e port==443 -o target
    \b
    # 输出端口为443,title中包含系统的url到命令行
    python format.py input.json -e port==443 -e title::系统 -o url
    \b
    # 与上相同,结果输出到res.txt文件
    python format.py input.json -e port==443 -e title::系统 -o target -f res.txt


    """
    if not outfile:
        outfunc = print
    else:
        outfunc = partial(toFile, outfile)

    results = []
    for file in files:
        print(file.name)
        results += loadResult(file)

    results = GogoResults(results)
    results = results.exprs(exprs, _or)
    if output == "full":
        outfunc(results.full)
    elif output == "json":  # 输出过滤后的json
        outfunc(results.raw_json)
    elif output == "zombie":  # 输出结果到zombie
        outfunc(results.zombie)
    else:
        outputs = output.split(",")
        contents = results.getstrs(*outputs)
        outfunc("\n".join(contents))


if __name__ == '__main__':
    main()

