package comments

import "testing"

func TestCreateCommentID(t *testing.T) {

	id, _, err := createCommentId(
		"The repeating is annoying and unhelpful...",
		"4bbdfd53f295d3ac16fe003df113f0568cc81b93")
	if err != nil {
		t.Error(err)
	}

	expected := "bd6b3c881e48f38fb689815a6d4cc36e285db16c68fa88b9f7ca7824d3a4c5f7"
	if id != expected {
		t.Errorf("id generated does not match examples from production\n expected %s\n got      %s", expected, id)
	}

	id, _, err = createCommentId(
		`Probably the same for everyone, since this is an MP3 audio file :-)`,
		"36b7bd81c1f975878da8cfe2960ed819a1c85bb5")
	if err != nil {
		t.Error(err)
	}
	expected = "ffffeec9fd1f02216a18db6b73864da3b2cfb8660507f7c36046f3490abbce71"
	if id != expected {
		t.Errorf("id generated does not match examples from production\n expected %s\n got      %s", expected, id)
	}
}
