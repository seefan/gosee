package utils

import (
	"math"
)

//点数据源
type Point struct {
	X, Y float64
}

// 功能：判断点是否在多边形内，简单算法，如需要精确，还可以用弧长法验证
// 方法：求解通过该点的水平线与多边形各边的交点
// 结论：单边交点为奇数，成立!
//参数：
// p 待判断的某个点
// lp  多边形的各个顶点坐标（首末点可以不一致）
func PointInsidePolygon1(p *Point, lp []*Point) bool {
	nCount := len(lp)
	nCross := 0
	for i := 0; i < nCount; i++ {
		p1 := lp[i]
		p2 := lp[(i+1)%nCount]
		// 求解 Y=p.Y 与 p1p2 的交点
		if p1.Y == p2.Y { // p1p2 与 Y=p0.y平行
			continue
		}
		if p.Y < math.Min(p1.Y, p2.Y) { // 交点在p1p2延长线上
			continue
		}
		if p.Y >= math.Max(p1.Y, p2.Y) {
			// 交点在p1p2延长线上
			continue
		}
		// 求交点的 x 坐标 --------------------------------------------------------------
		x := (p.Y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y) + p1.X
		if x > p.X {
			nCross++ // 只统计单边交点
		}
	}
	// 单边交点为偶数，点在多边形之外 ---
	return (nCross%2 == 1)
}
func get_tmp(p0 *Point) int {
	if p0.X >= 0 {
		if p0.Y >= 0 {
			return 0
		} else {
			return 3
		}
	} else {
		if p0.Y >= 0 {
			return 1
		} else {
			return 2
		}
	}
}

// 功能：判断点是否在多边形内，简单算法，如需要精确
// 方法：弧长法验证
// 参数：
// p 待判断的某个点
// lp  多边形的各个顶点坐标（首末点可以不一致）
func PointInsidePolygon2(p *Point, lp []*Point) bool {
	if len(lp) == 0 {
		return false
	}
	var i, t1, t2, sum int

	n := len(lp)
	for i := 0; i <= n; i++ {
		lp[i].X -= p.X
		lp[i].Y -= p.Y
	}
	t1 = get_tmp(lp[0])

	for i = 1; i <= n; i++ {
		if lp[i].X == 0 && lp[i].Y == 0 {
			break
		}
		f := lp[i].Y*lp[i-1].X - lp[i].X*lp[i-1].Y
		if f > 0 && lp[i-1].X*lp[i].X <= 0 && lp[i-1].Y*lp[i].Y <= 0 {
			break
		}
		t2 = get_tmp(lp[i])
		if t2 == (t1+1)%4 {
			sum += 1
		} else if t2 == (t1+3)%4 {
			sum -= 1
		} else if t2 == (t1+2)%4 {
			if f > 0 {
				sum += 2
			} else {
				sum -= 2
			}
		}
		t1 = t2
	}
	if i <= n || sum > 0 {
		return true
	}
	return false
}
