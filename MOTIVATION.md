# Introduction

I'm by no means a good developer, nor engineer, mostly due to the fact that I have no experience on "real" software projects: most of what I've built have been small-ish projects built at most by two people, with no real scale. However, for quite some time now, I've been finding myself making concessions constantly on my selection of tech stacks, and on coding practices themselves (the latter in most cases following from the former, but it'd be unfair of me to not take the blame where balme is due).

By concessions, I mean especifically choosing "hype" libraries/framework, just because they were the New Big Thing, promissing espectacular developer experience. I find that they've always failed at that. At most, they made it faster for me to reach a stage where I was actually "building and shipping, bro" (false, I never shipped anything meaningful), but was that worth it if I never really enjoyed the process? Not only that, but the moment I tried to deviate slightly from the way the developers of said tools wanted me to do things (out of necessity or out of stubbornness) I found myself waisting hours trying to fight a tool that was supposed to help me do things faster. In most cases, by the end of that struggle I generally dropped the project (the energy having been zapped from me), or decided that I'd made a bad decision in choosing whatever it is I choosed, and had, thus, to start the project again.

 This project is an attempt at fixing that. I intend to implement this backend template doing things the way I like, forcing me to avoid taking the easy road. After careful thought, I've decided this implies:
 1. As little dependencies as possible. I'm using go, claimed to have a really complete stdlib, so I should strive to implement everything possible by hand. I'll be rebuilding a lot of wheels, except in cases where that is either not possible, or straight up dumb. Other than that, I'll strive to use the stdlib and my brain as much as possible. This means no web frameworks, no ORMs, probably no validator libraries (I'll see).
 2. No AI-generated code. This is both a preparation for future work and, most importantly, a learning experience. AI has little to no use in the latter. I might ask conceptual questions if I'm ever stuck on something, but I won't ask it to generate code for me.
 3. No premature optimization. I've taken some decisions before based on questions like "what happens if I get 1M uses?" Truth be told, I think I could count the users of the systems I've developed in a single hand. Until I have to deal with scaling issues, I won't optimize for them. Code is maleable and cheap, if I need to change something at a later date because I need to, then I will, but I'll try to waste no time in hypotheticals.
 4. Code structure: I lost some time thinking about this. Should I go with a clean architecture, to make it easy to change technologies later? Should I use a vertical slice architecture to make it easies to implement new features? I've come to the conclussion (actually [Casey Muratori and ThePrimagen)[https://www.youtube.com/watch?v=DsAclZbP_Us] made me reach it) that this is nothing more than premature optimization. This is a template that I intend to use in multiple projects that will have multiple requirements. If I ever need a pre-defined architecture, I'll refactor things. For now, I'll follow the architecture I know:
 ```
src/
├─ cmd/
│  ├─ api/
│  │  ├─ v1.go
│  ├─ main.go
├─ models/
│  ├─ user.go
├─ utils/
│  ├─ db.go
│  ├─ config.go
├─ sevices/
│  ├─ auth/
│  │  ├─ routes.go
│  │  ├─ repository.go
 ```
 There are probably a lot of things to improve on this. I'm open to all critiques. I am, however, going to move forward with this.
 5. Testing. I've never managed to settle on a testing framework that works for me. I hate mocking things, and I like TDD only as an exploratory tool when I'm not sure how to build something (even then, it's more a matter of building some tests to see if requirements are met, does testing that a summing function returns 2+2=4, 2+0=2 and so on and so forth really give you meaningful information?). The testing "framework" I've reached is using pytest for "integration" tests (i.e. testing that endpoints return what they should return. You could call these system tests I guess, but I don't care much for semantics) and using go for unit tests. I'll be also following ThePrimeagen advice here, since it's the one that most closely ressembles what I believe (this time you'll have to find the video on your own, tho): if I want to test something, I'll pull that out in a function and test that specifically. If I have to mock something, then I'm probably doing something wrong.

With that out of the way, let's get building

# 1. Configuring Things

I believe that the first step should always be to get something working. As such, I'll implement a simple endpoint that connects to a database and stores/retrieves things from it. This will be the stepping stone on which I'll later build the auth system.
