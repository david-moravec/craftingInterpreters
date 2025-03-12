import os
import argparse
import subprocess
from dataclasses import dataclass


def recompile_golox():
    os.chdir("../golox")
    subprocess.run(["go", "build", "."])
    os.chdir("../test")


def recompile_clox():
    pass


def run_clox_exe():
    return


def run_golox_exe(filepath: os.PathLike):
    print(f"Running: {filepath}")
    test = Test(filepath)
    return test.run()


def test_lox_file_golox(filepath: os.PathLike):
    return run_golox_exe(filepath)


def test_lox_file_clox(filepath: os.PathLike):
    print(f"clox test {filepath}")
    return


def benchmark_lox(test_function):
    for file in os.listdir("benchmark"):
        test_function("benchmark" + os.sep + file)


@dataclass(frozen=True, kw_only=True)
class LineInfo:
    no: int
    code: str
    expected: str

    def __str__(self) -> str:
        return f"[Line: {'{line_no:2d}'.format(line_no=self.no)}]: {self.code}\n[Expected]: {self.expected.decode('ascii')}"


class Test:
    def __init__(self, path: os.PathLike):
        self._path = path
        self._expected = []
        self.parse()

    def run(self):
        p = subprocess.Popen(
            ["../golox/golox.exe", self._path],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )

        return self.test_output(p.stdout)

    def test_output(self, out_buff):
        result = 0
        for line_info, output in zip(self._expected, out_buff):
            output = output.rstrip()
            expected: str = line_info.expected
            if expected != output:
                try:
                    expected = float(expected)
                    output = float(output)
                except ValueError:
                    try:
                        expected = int(expected)
                        output = int(output)
                    except ValueError:
                        pass

            if expected != output:
                result = 1
                output = (
                    output.decode("ascii")
                    if isinstance(output, bytes)
                    else output
                )
                print("\n############## Error ###################")
                print(line_info)
                print(f"[Got]     : {output}")
                print("########################################\n")

        return result

    def parse(self):
        with open(self._path) as file:
            for i, line in enumerate(file):
                split = line.split(sep="// expect: ")

                if len(split) == 2:
                    expected: str = split[1]
                    expected = expected.rstrip()
                    expected = expected.encode("ascii")
                    self._expected.append(
                        LineInfo(
                            no=i, code=split[0].lstrip(), expected=expected
                        )
                    )


def run_tests(args: argparse.Namespace):
    if args.golox:
        test_function = test_lox_file_golox
        recompile_function = recompile_golox
    else:
        test_function = test_lox_file_clox
        recompile_function = recompile_clox

    if args.recompile:
        recompile_function()

    if args.benchmark:
        benchmark_lox(test_function)
    else:
        root = args.path if args.path else "./"

        if not os.path.isdir(root):
            test_function(root)
        else:
            failed = 0
            total = 0
            for subdir, dirs, files in os.walk(root):
                if subdir in ("./benchmark", "./expressions", "./scanning"):
                    continue

                for file in files:
                    filepath = subdir + os.sep + file

                    if filepath.endswith(".lox"):
                        total += 1
                        failed += test_function(filepath)

            print(f"Failed: {failed} ({failed/total:.0%})\n ")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(prog="Test")
    parser.add_argument("-r", "--recompile", action="store_true")
    parser.add_argument("-g", "--golox", action="store_true")
    parser.add_argument("-b", "--benchmark", action="store_true")
    parser.add_argument("-p", "--path")
    run_tests(parser.parse_args())
