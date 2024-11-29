package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func Test_getWordStats(t *testing.T) {
	tests := []struct {
		name string
		text string
		want map[string]int
	}{
		{
			name: "all words are different",
			text: "жила-была алиса в стране чудес",
			want: map[string]int{
				"жила-была": 1,
				"алиса":     1,
				"в":         1,
				"стране":    1,
				"чудес":     1,
			},
		},
		{
			name: "using non-similar words",
			text: "тест Тест тест, тест тест тест, Тест, тест Тест тест2",
			want: map[string]int{
				"тест":  9,
				"тест2": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getWordStats(tt.text))
		})
	}
}

func Test_topKeys(t *testing.T) {
	type args struct {
		m     map[string]int
		count int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "all words with same count",
			args: args{
				m: map[string]int{
					"жила-была": 1,
					"алиса":     1,
					"в":         1,
					"стране":    1,
					"чудес":     1,
				},
				count: 5,
			},
			want: []string{
				"алиса",
				"в",
				"жила-была",
				"стране",
				"чудес",
			},
		},
		{
			name: "some words with same count, some words with different",
			args: args{
				m: map[string]int{
					"аист":  2,
					"чайка": 3,
					"кошка": 3,
					"робот": 1,
					"целый": 1,
				},
				count: 5,
			},
			want: []string{
				"кошка",
				"чайка",
				"аист",
				"робот",
				"целый",
			},
		},
		{
			name: "less words than map length",
			args: args{
				m: map[string]int{
					"аист":  2,
					"чайка": 3,
					"кошка": 3,
					"робот": 1,
					"целый": 1,
				},
				count: 2,
			},
			want: []string{
				"кошка",
				"чайка",
			},
		},
		{
			name: "count more than map length",
			args: args{
				m: map[string]int{
					"аист":  2,
					"чайка": 3,
					"кошка": 3,
					"робот": 1,
					"целый": 1,
				},
				count: 10,
			},
			want: []string{
				"кошка",
				"чайка",
				"аист",
				"робот",
				"целый",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, topKeys(tt.args.m, tt.args.count))
		})
	}
}

func Test_normalizeWord(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "word without punctuation marks",
			args: args{
				word: "hello",
			},
			want: "hello",
		},
		{
			name: "word with commas",
			args: args{
				word: ",hello,,,",
			},
			want: "hello",
		},
		{
			name: "word with dot and commas",
			args: args{
				word: ".hello,,,",
			},
			want: "hello",
		},
		{
			name: "word is -",
			args: args{
				word: "-",
			},
			want: "",
		},
		{
			name: "word is ----",
			args: args{
				word: "----",
			},
			want: "----",
		},
		{
			name: "word with middle - ",
			args: args{
				word: ",one-two..-",
			},
			want: "one-two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, normalizeWord(tt.args.word), "normalizeWord(%v)", tt.args.word)
		})
	}
}
