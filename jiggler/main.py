import pyautogui as pag
import time
import argparse
from random import randint

class Jiggler():
    def __init__(self):
        ap = argparse.ArgumentParser(formatter_class=argparse.MetavarTypeHelpFormatter)
        ap.add_argument('--loops', type=int, default=3)
        ap.add_argument('--sleep', type=float, default=2)
        ap.add_argument('--move-by', type=int, default=4)
        ap.add_argument('--random', action='store_true', default=False)
        ns = ap.parse_args()

        self.LOOPS = ns.loops
        self.MODIFIER = ns.move_by
        self.SLEEP = ns.sleep
        self.RANDOM = ns.random
        pag.FAILSAFE = False

    def up(self) -> None:
        pag.moveRel(0, -self.MODIFIER)
        pass

    def down(self) -> None:
        pag.moveRel(0, self.MODIFIER)
        pass

    def left(self) -> None:
        pag.moveRel(-self.MODIFIER, 0)
        pass

    def right(self) -> None:
        pag.moveRel(self.MODIFIER, 0)
        pass

    def jiggle(self) -> None:
        self.backup = self.MODIFIER
        for i in range(0, self.LOOPS):
            if self.RANDOM:
                self.MODIFIER = randint(1, 64)
            else:
                self.MODIFIER = self.MODIFIER * (i+1)
            self.right()
            self.down()
            self.left()
            self.up()
        self.MODIFIER = self.backup

if __name__ == '__main__':
    j = Jiggler()
    while True:
        j.jiggle()
        time.sleep(j.SLEEP)