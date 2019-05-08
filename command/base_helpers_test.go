package command

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestParseArgsData(t *testing.T) {
	t.Parallel()

	t.Run("stdin_full", func(t *testing.T) {
		t.Parallel()

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(`{"foo":"bar"}`))
			stdinW.Close()
		}()

		m, err := parseArgsData(stdinR, []string{"-"})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("stdin_value", func(t *testing.T) {
		t.Parallel()

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(`bar`))
			stdinW.Close()
		}()

		m, err := parseArgsData(stdinR, []string{"foo=-"})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})
	t.Run("file_full", func(t *testing.T) {
		t.Parallel()

		f, err := ioutil.TempFile("", "vault")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte(`{"foo":"bar"}`))
		f.Close()
		defer os.Remove(f.Name())

		m, err := parseArgsData(os.Stdin, []string{"@" + f.Name()})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("file_value", func(t *testing.T) {
		t.Parallel()

		f, err := ioutil.TempFile("", "vault")
		if err != nil {
			t.Fatal(err)
		}
		f.Write([]byte(`bar`))
		f.Close()
		defer os.Remove(f.Name())

		m, err := parseArgsData(os.Stdin, []string{"foo=@" + f.Name()})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "bar" {
			t.Errorf("expected %q to be %q", v, "bar")
		}
	})

	t.Run("file_value_escaped", func(t *testing.T) {
		t.Parallel()

		m, err := parseArgsData(os.Stdin, []string{`foo=\@`})
		if err != nil {
			t.Fatal(err)
		}

		if v, ok := m["foo"]; !ok || v != "@" {
			t.Errorf("expected %q to be %q", v, "@")
		}
	})
}

func TestUsernameFromEnv(t *testing.T) {
	userTests := []struct {
		name string
		env  string
		want string
	}{
		{name: "Test LOGNAME", env: "LOGNAME", want: "test"},
		{name: "Test USER", env: "USER", want: "test"},
		{name: "Test empty", env: "", want: ""},
	}
	for _, tt := range userTests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("LOGNAME")
			os.Unsetenv("USER")
			os.Setenv(tt.env, tt.want)
			got := UsernameFromEnv()
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPasswordFromEnv(t *testing.T) {
	passwordTests := []struct {
		name string
		env  string
		want string
	}{
		{name: "Test PASSWORD", env: "PASSWORD", want: "test"},
		{name: "Test empty password", env: "", want: ""},
	}
	for _, tt := range passwordTests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("PASSWORD")
			os.Setenv(tt.env, tt.want)
			got := PasswordFromEnv()
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestEnsureTrailingSlash(t *testing.T) {
	got := EnsureTrailingSlash("test")
	want := "test/"
	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}
}
