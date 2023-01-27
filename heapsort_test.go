package qsort_test

// func Test_qsort(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input []string
// 	}{
// 		{
// 			name:  "len <= 7",
// 			input: lmacho.CpuNames()[:6],
// 		},
// 		{
// 			name:  "7 < len < 40",
// 			input: lmacho.CpuNames()[:9],
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := make([]lmacho.FatArch, len(tt.input))
// 			want := make([]lmacho.FatArch, len(got))
// 			copy(want, got)

// 			qsort.Slice(got, cmpFunc)
// 			if !reflect.DeepEqual(want, got) {
// 				t.Errorf("\nwant: %v\ngot: %v", cnv(want), cnv(got))
// 			}
// 		})
// 	}
// }
