# goenvsubst

goenvsubst (_Go Environment substitution_; gopherified version of [envsubst](https://www.gnu.org/software/gettext/manual/html_node/envsubst-Invocation.html))
provides a simple way to write environment variable values to files on container startup.

Let's say you have nginx in a container and simply want to dynamically set a value from the container environment to its configuration,
`goenvsubst` is an easy way to do that. 

It supports
 
* Replace ENV Variable values, with optional default values (simple `${ENVNAME}` or `${ENVNAME:-default}` with default value)
* Full go templating support 

## Basic usage

By default (with no arguments), `goenvsubst` reads input from _stdin_:

```
echo 'This is my Template with ${VAR}:' > tmpl
export VAR="my dynamic value"
cat tmpl | goenvsubst > config.conf

cat config.conf
This is my Template with my dynamic value:
```

You can set default values inline. If the ENV variable is not set, that value will be set:

```
echo 'My value is ${VAR:-default} set.' > tmpl
cat tmpl | goenvsubst > config.conf

cat config.conf
My value is default set.
```

You can also pass the template file path as the first argument:

```
echo 'This is my Template with ${VAR}:' > tmpl
export VAR="my dynamic value"
goenvsubst ./tmpl > config.conf

cat config.conf
This is my Template with my dynamic value:
```

### Go templating

Before replacing any `${..}` occurrences, the whole content is parsed as a Go template. Every ENV variable
is a variable in the template.

```
echo 'My name is {{.NAME}} {{ if .AGE }}and i am {{.AGE}} years old{{ end }} :' > tmpl
export NAME="Jeff"
export AGE=13
goenvsubst ./tmpl > config.conf

cat config.conf
My name is Jeff and i am 13 years old :
```

## Use in your container

Just add the latest [prebuilt binary for 64-bit Linux][releases] from the Github releases page to your container:

Dockerfile
```
FROM ....

ENV GOENVSUBST_VERSION=v0.1.0
ADD https://github.com/libgraviton/goenvsubst/releases/download/${GOENVSUBST_VERSION}/goenvsubst-amd64 /sbin/goenvsubst

RUN chmod +x /sbin/goenvsubst
```

## Notes

I did this for my own need. After doing it, I discovered [a8m/envsubst](https://github.com/a8m/envsubst) which seems
to provide a more mature set of checking features and more advanced syntax.

I contemplated first to delete this again. But then I realized that this small utility is a better match for many simple use cases
and as it also supports Go templating (which a8m/envsubst does not), I decided to keep this. 