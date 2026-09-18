package database

import "strings"

// LCSSimilarity returns the normalized longest-common-subsequence score in [0,1].
// Comparison is case-insensitive and ignores surrounding whitespace.
func LCSSimilarity(left, right string) float64 { a,b:=[]rune(strings.ToLower(strings.TrimSpace(left))),[]rune(strings.ToLower(strings.TrimSpace(right)));if len(a)==0&&len(b)==0{return 1};if len(a)==0||len(b)==0{return 0};previous:=make([]int,len(b)+1);for _,ra:=range a{current:=make([]int,len(b)+1);for j,rb:=range b{if ra==rb{current[j+1]=previous[j]+1}else if current[j]>previous[j+1]{current[j+1]=current[j]}else{current[j+1]=previous[j+1]}};previous=current};score:=previous[len(b]);if len(a)>len(b){return float64(score)/float64(len(a))};return float64(score)/float64(len(b))}
func SimilarAtLeast(left,right string, threshold float64) bool { return LCSSimilarity(left,right)>=threshold }
const similarityThreshold=0.75
func containsSimilar(query string, values ...string) bool { for _,value:=range values{if SimilarAtLeast(query,value,similarityThreshold){return true}};return false }
