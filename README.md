# goerror

Its is a super-small super-light struct that I use in most of my projects to wrap errors.
It is a little bit unorthodox, since it features some chaining methods, but I think it is very handy.

Each error has a unique ID, and a message. The ID is used to identify the error, e.g. with errors.Is().
The message is used to display a user message.

This originated from one of my web projects. Error.Message is thought as return value for the user, the actual error is stored in Error.Err.
Therefore, Error.MarshalJSON() will return the message, and Error.Id, but not the actual error.