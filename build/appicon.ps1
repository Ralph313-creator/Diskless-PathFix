# Draws build\appicon.png (1024x1024), the RapTech Portable app icon: a white
# folder (the app's files) with a badge holding an outward arrow (portable),
# on a teal-to-blue rounded square. wails build makes build\windows\icon.ico
# from appicon.png when icon.ico is missing, so delete icon.ico after a change.
Add-Type -AssemblyName System.Drawing
$size = 1024
$bmp = New-Object System.Drawing.Bitmap $size, $size
$g = [System.Drawing.Graphics]::FromImage($bmp)
$g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
$g.Clear([System.Drawing.Color]::Transparent)

function RoundRect([float]$x, [float]$y, [float]$w, [float]$h, [float]$r) {
    $p = New-Object System.Drawing.Drawing2D.GraphicsPath
    $d = 2 * $r
    $p.AddArc($x, $y, $d, $d, 180, 90)
    $p.AddArc($x + $w - $d, $y, $d, $d, 270, 90)
    $p.AddArc($x + $w - $d, $y + $h - $d, $d, $d, 0, 90)
    $p.AddArc($x, $y + $h - $d, $d, $d, 90, 90)
    $p.CloseFigure()
    return $p
}

# Background
$bg = RoundRect 32 32 960 960 200
$grad = New-Object System.Drawing.Drawing2D.LinearGradientBrush (New-Object System.Drawing.Point 32, 32), (New-Object System.Drawing.Point 992, 992), ([System.Drawing.Color]::FromArgb(255, 20, 184, 166)), ([System.Drawing.Color]::FromArgb(255, 29, 78, 216))
$g.FillPath($grad, $bg)

# Folder: tab, then body
$white = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::White)
$g.FillPath($white, (RoundRect 176 250 300 160 48))
$g.FillRectangle($white, 176, 300, 300, 150) # tab joins the body without a notch
$g.FillPath($white, (RoundRect 176 330 672 470 64))

# Badge with an arrow pointing up and out
$dark = New-Object System.Drawing.SolidBrush ([System.Drawing.Color]::FromArgb(255, 15, 23, 42))
$g.FillEllipse($dark, 556, 536, 360, 360)
$pen = New-Object System.Drawing.Pen ([System.Drawing.Color]::White), 52
$pen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
$pen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round
$pen.LineJoin = [System.Drawing.Drawing2D.LineJoin]::Round
$g.DrawLine($pen, 660, 792, 820, 632)
$g.DrawLines($pen, [System.Drawing.Point[]]@((New-Object System.Drawing.Point 690, 632), (New-Object System.Drawing.Point 820, 632), (New-Object System.Drawing.Point 820, 762)))

$g.Dispose()
$out = Join-Path $PSScriptRoot 'appicon.png'
$bmp.Save($out, [System.Drawing.Imaging.ImageFormat]::Png)
$bmp.Dispose()
"wrote $out"
