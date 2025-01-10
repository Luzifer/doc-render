\documentclass[10pt,a4paper]{article}
\usepackage{graphicx}
\usepackage{lastpage}
\usepackage{anysize}
\usepackage{fancyhdr}
\usepackage[ngerman,shorthands=off]{babel}
\usepackage[T1]{fontenc}
\usepackage[utf8]{inputenc}
\usepackage{aeguill}
\usepackage{eurosym}
\usepackage{multicol}

\usepackage{fontspec}
\setmainfont[
  Path           = fonts/ ,
  Extension      = .ttf ,
  UprightFont    = *-Regular ,
  BoldFont       = *-Bold ,
  BoldItalicFont = *-BoldItalic ,
  ItalicFont     = *-Italic
]{Roboto}

\setmonofont[
  Path           = fonts/ ,
  Extension      = .ttf ,
  UprightFont    = *-Regular ,
  BoldFont       = *-Bold ,
  BoldItalicFont = *-BoldItalic ,
  ItalicFont     = *-Italic
]{RobotoCondensed}

%%% Format der Seite anpassen
\marginsize{2.5cm}{2.5cm}{0.85cm}{3.41cm}

\def\code#1{\texttt{#1}}

%%% Header und Footer bauen
\pagestyle{fancy}
\renewcommand{\footrulewidth}{0.4pt}
\renewcommand{\headrulewidth}{0pt}
\renewcommand{\headheight}{2.75cm}
\renewcommand{\headsep}{0cm}
\setlength{\unitlength}{1mm}
\lhead{\Huge Max Muster \Large\\ \small\ \\}
\chead{}
\rhead{\large{Musterstr. 123\\12345 Musterstadt\\\tiny{\ }\\}}
\lfoot{\small{max@muster.io\\
  +49 123 4567890}}
\cfoot{\thepage\ / \pageref{LastPage}}
\rfoot{\small{example.com/web\\
  example.com/github}}

\begin{document}
\setlength{\parindent}{0mm}
\setlength{\parskip}{6pt}

%%% Dokumentenkrams nur auf der ersten Seite

\vspace{8mm}
\hfill

\begin{picture}(0,0)
  \put( -4, -3.85){\makebox(85,4){\fontsize{2.5mm}{2.5mm}\selectfont{Max Muster\ $\cdot$\ Musterstr. 123\ $\cdot$\ 12345 Musterstadt}}}
  \put( -4, -3.95){\line(1,0){85}}
  \put(3,-15){\parbox[t]{3in}{
    {{ md2tex .Values.address }}
  }}
  \put(-4,0){\line( 1, 0){1}} \put(-4,0){\line( 0,-1){1}}
  \put(81,0){\line(-1, 0){1}} \put(81,0){\line( 0,-1){1}}
  \put(-4,-40.85){\line( 1, 0){1}} \put(-4,-40.85){\line( 0, 1){1}}
  \put(81,-40.85){\line(-1, 0){1}} \put(81,-40.85){\line( 0, 1){1}}
\end{picture}
\hfill

\vspace{5.0cm}

\makebox[40.0mm][l]{}
\makebox[40.0mm][l]{}
\makebox[40.0mm][l]{}
\makebox[36.0mm][r]{Datum} \\
\makebox[40.0mm][l]{}
\makebox[40.0mm][l]{}
\makebox[40.0mm][l]{}
\makebox[36.0mm][r]{ {{- default "\\today" .Values.date -}} }

\vspace{1cm}

%%% Dokumenteninhalt

\textbf{ {{- .Values.subject -}} }

\vspace{0.5cm}

{{ md2tex .Values.text }}
\end{document}
