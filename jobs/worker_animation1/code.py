#python
from manim import *

class ImpedanceTriangle(Scene):
    def construct(self):
        # Configuration
        R_len = 4
        X_len = 3
        origin = LEFT * 2 + DOWN * 1.5

        # Vectors
        vec_R = Arrow(origin, origin + RIGHT * R_len, buff=0, color=BLUE)
        vec_X = Arrow(vec_R.get_end(), vec_R.get_end() + UP * X_len, buff=0, color=RED)
        vec_Z = Arrow(origin, vec_X.get_end(), buff=0, color=YELLOW)

        # Labels
        label_R = MathTex("R").next_to(vec_R, DOWN)
        label_X = MathTex("X").next_to(vec_X, RIGHT)
        label_Z = MathTex("Z").move_to(vec_Z.get_center() + UP * 0.5 + LEFT * 0.5)

        # Angles
        right_angle = RightAngle(vec_R, vec_X, length=0.4, quadrant=(-1, 1))
        phi = Angle(vec_R, vec_Z, radius=0.7)
        label_phi = MathTex(r"\phi").next_to(phi, RIGHT).shift(UP * 0.1)

        # Title and Formulas
        title = Text("Impedance Triangle").to_edge(UP)
        formulas = VGroup(
            MathTex(r"Z = \sqrt{R^2 + X^2}"),
            MathTex(r"\tan \phi = \frac{X}{R}")
        ).arrange(DOWN, aligned_edge=LEFT).to_corner(UL).shift(DOWN)

        # Animation Sequence
        self.play(Write(title))
        self.play(GrowArrow(vec_R), Write(label_R))
        self.play(GrowArrow(vec_X), Write(label_X))
        self.play(Create(right_angle))
        self.play(GrowArrow(vec_Z), Write(label_Z))
        self.play(Create(phi), Write(label_phi))
        self.play(Write(formulas))
        self.wait(2)
