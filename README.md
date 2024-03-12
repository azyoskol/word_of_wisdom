## Task
Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge

## Explanations

### POW algorithm

I did research and I think it's my favorite choice [Memory-bound](https://en.wikipedia.org/wiki/Memory-bound_function). 

Why do I choice this one? The general requarement is protected from DDOS attacks with Proof of Work.
As you could see in article Memory-bound function was create to prevent email spamming. The idea MBound functions would consume CPU resources at the client's machine for each request, thus preventing huge amounts of requests from being sent in a short period. In addition, we increase the cost of client resources, thus DDOS attacks will be very expensive.
