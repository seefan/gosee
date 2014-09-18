package utils

import (
	"log"
	"math"
	"math/rand"
)

const (
	EPS = 1e-8
)

//点数据源
type Point struct {
	X, Y float64
}

//一个矩形
type Box struct {
	MinX, MinY, MaxX, MaxY float64
}

//求取两点之间的距离
func GetDistance(tpt1, tpt2 *Point) float64 {
	x := math.Abs(tpt1.X - tpt2.X)
	y := math.Abs(tpt1.Y - tpt2.Y)
	return math.Sqrt(x*x + y*y)
}

//判断点是否在线上
func JudgePtInLine(tpt1, tpt2, tpt *Point) bool {
	dx1 := GetDistance(tpt1, tpt2)
	dx2 := GetDistance(tpt, tpt1)
	dx3 := GetDistance(tpt, tpt2)
	dx := dx3 + dx2 - dx1
	if dx >= -0.0000000001 && dx <= 0.0000000001 {
		return true
	}
	return false
}

// 功能：判断点是否在多边形内
// 方法：求解通过该点的水平线与多边形各边的交点
// 结论：单边交点为奇数，成立!

//参数：
// POINT p 指定的某个点
// LPPOINT ptPolygon 多边形的各个顶点坐标（首末点可以不一致）
// int nCount 多边形定点的个数

func PtInPolygon(p *Point, ptPolygon []*Point) bool {
	nCross := 0
	nCount := len(ptPolygon)
	for i := 0; i < nCount; i++ {
		p1 := ptPolygon[i]
		p2 := ptPolygon[(i+1)%nCount]
		// 求解 y=p.Y 与 p1p2 的交点
		if p1.Y == p2.Y { // p1p2 与 y=p0.y平行
			continue
		}
		if p.Y < math.Min(p1.Y, p2.Y) { // 交点在p1p2延长线上
			continue
		}
		if p.Y >= math.Max(p1.Y, p2.Y) { // 交点在p1p2延长线上
			continue
		}
		// 求交点的 X 坐标 --------------------------------------------------------------
		x := (p.Y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y) + p1.X

		if x > p.X {
			nCross++ // 只统计单边交点
		}
	}
	// 单边交点为偶数，点在多边形之外
	return (nCross%2 == 1)
}

//取一组点的最大范围
func getBoundingBox(poly []*Point) (box *Box) {
	n := len(poly)
	box = new(Box)
	if n == 0 {
		return
	}
	box.MinX = poly[0].X
	box.MaxX = poly[0].X
	box.MinY = poly[0].Y
	box.MaxY = poly[0].Y
	for i := 1; i < n; i++ {
		if poly[i].X > box.MaxX {
			box.MaxX = poly[i].X
		}
		if poly[i].X < box.MinX {
			box.MinX = poly[i].X
		}
		if poly[i].Y > box.MaxY {
			box.MaxY = poly[i].Y
		}
		if poly[i].Y < box.MinY {
			box.MinY = poly[i].Y
		}
	}
	return
}

//判断点在Bounding Box以内还是以外。
func pOutsideBoundingBox(box *Box, p *Point) bool {
	if p.X < box.MinX || p.X > box.MaxX {
		return true
	}
	if p.Y < box.MinY || p.Y > box.MaxY {
		return true
	}
	return false
}

//计算３点构成的向量的叉乘积，叉乘积为０则３点一直线。
func crossPro(p1, p2, p *Point) float64 {
	return ((p1.X-p.X)*(p2.Y-p.Y) - (p2.X-p.X)*(p1.Y-p.Y))
}

//判断点是否在边线上或顶点上。
func pOnline(p1, p2, p *Point) bool {
	if (p.X == p1.X) && (p.Y == p1.Y) {
		return true
	}
	if (p.X == p2.X) && (p.Y == p2.Y) {
		return true
	}
	if (p.X >= p1.X) && (p.X >= p2.X) {
		return false
	}
	if (p.X <= p1.X) && (p.X <= p2.X) {
		return false
	}
	if (p.Y >= p1.Y) && (p.Y >= p2.Y) {
		return false
	}
	if (p.Y <= p1.Y) && (p.Y <= p2.Y) {
		return false
	}
	area2 := crossPro(p1, p2, p)
	if math.Abs(area2) < EPS {
		return true
	}
	return false
}

// if (t2==0) AB & CD are parallel
// if (t2==0 && t1==0) AB & CD are collinear.
// If 0<=r<=1 && 0<=s<=1, intersection exists
// r<0 or r>1 or s<0 or s>1 line segments do not intersect
// return 1 -- intersect
// return 0 -- not intersect
// return 0 --  (t2==0  && t1 != 0) parallel
// return -1 -- ( t2==0 && t1==0)  collinear.
//计算线段AB和CD相交状态和交点。
func intersect(A, B, C, D *Point) int {
	var r, s float64
	var t1, t2, t3 float64
	t2 = (B.X-A.X)*(D.Y-C.Y) - (B.Y-A.Y)*(D.X-C.X)
	t1 = (A.Y-C.Y)*(D.X-C.X) - (A.X-C.X)*(D.Y-C.Y)
	if (math.Abs(t2) < EPS) && (math.Abs(t1) < EPS) {
		return -1
	}
	if math.Abs(t2) < EPS && math.Abs(t1) >= EPS {
		return 0
	}
	r = t1 / t2
	if (r < 0.0) || (r > 1.0) {
		return 0
	}
	t3 = (A.Y-C.Y)*(B.X-A.X) - (A.X-C.X)*(B.Y-A.Y)
	s = t3 / t2
	if (s < 0.0) || (s > 1.0) {
		return 0
	}
	return 1
}

//判断点在面内
func PointInsidePolygon2(p *Point, poly []*Point) bool {
	return PointInside(p, poly) != 0
}

//判断点在面内
// return 0 -- outside the polygon
// return 1 -- inside the polygon
// return 2 -- online the polygon
func PointInside(q *Point, poly []*Point) int {
	NN := len(poly)
	var flag_online, flag int
	NN_cross := 0
	qi := [9]Point{}
	qk := new(Point) // not use 0
	var i, k int
	var dy, dx, dy2, dx2 float64
	var shift_rd float64
	box := getBoundingBox(poly)
	//log.Println(box)
	//check 1 -- BoundingBox check
	if pOutsideBoundingBox(box, q) {
		return 0
	}
	//check 2 -- Vertex and Online check
	flag_online = 0
	for i := 0; i < NN-1; i++ {
		if pOnline(poly[i], poly[i+1], q) {
			flag_online = 1
			break
		}
	}

	if flag_online == 1 {
		return 2
	}
	// make qi[1]

	qi[1].X = box.MinX - 1.0
	qi[1].Y = q.Y
	qi[2].X = box.MaxX + 1.0
	qi[2].Y = q.Y
	qi[3].Y = box.MinY - 1.0
	qi[3].X = q.X
	qi[4].Y = box.MaxY + 1.0
	qi[4].X = q.X
	shift_rd = rand.Float64()
	qi[5].X = box.MinX - 1.0
	qi[5].Y = box.MinY - shift_rd
	qi[6].X = box.MaxX + shift_rd
	qi[6].Y = box.MinY - 1.0
	shift_rd = rand.Float64()
	qi[7].X = box.MaxX + 1.0
	qi[7].Y = box.MaxY + shift_rd
	qi[8].X = box.MinX - shift_rd
	qi[8].Y = box.MaxY + 1.0
	// which Ray
	k = 1
	flag = 0
	for i = 0; i < NN; i++ {
		if (poly[i].Y == q.Y) && (poly[i].X < q.X) {
			flag = 1
			break
		}
	}

	if flag == 0 {
		goto Lab_K
	}
	k = 2
	flag = 0
	for i = 0; i < NN; i++ {
		if (poly[i].Y == q.Y) && (poly[i].X > q.X) {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
	k = 3
	flag = 0

	for i = 0; i < NN; i++ {
		if (poly[i].X == q.X) && (poly[i].Y < q.Y) {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
	k = 4
	flag = 0
	for i = 0; i < NN; i++ {
		if (poly[i].X == q.X) && (poly[i].Y > q.Y) {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
	k = 5
	flag = 0
	dx = q.X - qi[5].X
	if math.Abs(dx) < EPS {
		goto Lab_6
	}
	dy = q.Y - qi[5].Y
	for i = 0; i < NN; i++ {
		dx2 = poly[i].X - qi[5].X
		dy2 = poly[i].Y - qi[5].Y
		if math.Abs(dy*dx2-dx*dy2) < EPS {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
Lab_6:
	k = 6
	flag = 0
	dx = q.X - qi[6].X
	if math.Abs(dx) < EPS {
		goto Lab_7
	}
	dy = q.Y - qi[6].Y
	for i = 0; i < NN; i++ {
		dx2 = poly[i].X - qi[6].X
		dy2 = poly[i].Y - qi[6].Y
		if math.Abs(dy*dx2-dx*dy2) < EPS {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
Lab_7:
	k = 7
	flag = 0
	dx = q.X - qi[7].X
	if math.Abs(dx) < EPS {
		goto Lab_8
	}
	dy = q.Y - qi[7].Y
	for i = 0; i < NN; i++ {
		dx2 = poly[i].X - qi[7].X
		dy2 = poly[i].Y - qi[7].Y
		if math.Abs(dy*dx2-dx*dy2) < EPS {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
Lab_8:
	;
	k = 8
	flag = 0
	dx = q.X - qi[8].X
	if math.Abs(dx) < EPS {
		goto Lab_9
	}
	dy = q.Y - qi[8].Y
	for i = 0; i < NN; i++ {
		dx2 = poly[i].X - qi[8].X
		dy2 = poly[i].Y - qi[8].Y
		if math.Abs(dy*dx2-dx*dy2) < EPS {
			flag = 1
			break
		}
	}
	if flag == 0 {
		goto Lab_K
	}
Lab_9:
	//printf("\007All Rays do not work -- use k=8 !\n");
	flag = 0
Lab_K:
	qk = &qi[k]
	//printf("Ray-%d\n",k);
	// check 3
	for i = 0; i < NN-1; i++ {
		flag = intersect(poly[i], poly[i+1], qk, q)
		if flag == 1 {
			NN_cross++
		}
	}
	if NN_cross%2 == 1 {
		// printf("point Q is inside of the polygon\n");
		return 1
	} else {
		// printf("point Q is outside of the polygon\n");
		return 0
	}
}
