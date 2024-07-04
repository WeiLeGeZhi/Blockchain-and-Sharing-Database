// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

contract Ballot {
    struct Voter {
        bool joined;
        bool voted;
    }

    struct Poll {
        string description;
        address creator;
        bool ended;
        address winner;
        address[] candidates;
        mapping(address => Voter) voters;
        mapping(address => uint) votes;
        uint voterCount;
    }

    struct SimplePoll {
        uint pollId;
        string description;
        address creator;
    }

    mapping(uint => Poll) public polls;
    uint public pollCount;
    address public owner;

    constructor() {
        owner = address(uint160(msg.sender));
    }

    event PollCreated(uint pollId, string description);
    event CandidateJoined(uint pollId, address candidate);
    event Voted(uint pollId, address voter, address candidate);
    event PollEnded(uint pollId, address winner);

    function createPoll(string memory description) public payable {
        pollCount++;
        Poll storage poll = polls[pollCount];
        poll.description = description;
        poll.creator = msg.sender;
        poll.ended = false;
        poll.voterCount = 0;
        poll.winner = address(0);

        emit PollCreated(pollCount, poll.description);
    }

    function joinPoll(uint pollId) public payable {
        Poll storage poll = polls[pollId];
        require(!poll.ended, "Poll has ended");
        require(!poll.voters[msg.sender].joined, "Already joined");
        require(!poll.voters[msg.sender].voted, "Already voted");
        require(polls[pollId].creator != address(0), "Poll does not exist");

        poll.voters[msg.sender].joined = true;
        poll.candidates.push(msg.sender);

        // payable(poll.creator).transfer(1);

        emit CandidateJoined(pollId, msg.sender);
    }

    function vote(uint pollId, address candidate) public {
        Poll storage poll = polls[pollId];
        require(!poll.ended, "Poll has ended");
        require(!poll.voters[msg.sender].voted, "Already voted");
        require(!poll.voters[msg.sender].joined, "Candidates cannot vote for themselves");
        require(poll.voters[candidate].joined, "Candidate not joined");
        require(polls[pollId].creator != address(0), "Poll does not exist");

        poll.voters[msg.sender].voted = true;
        poll.votes[candidate]++;
        poll.voterCount++;

        emit Voted(pollId, msg.sender, candidate);

        if (poll.voterCount >= 5) {
            endPoll(pollId);
        }
    }

    function endPoll(uint pollId) internal {
        Poll storage poll = polls[pollId];
        require(!poll.ended, "Poll has already ended");
        require(polls[pollId].creator != address(0), "Poll does not exist");

        address winningCandidate;
        uint maxVotes = 0;
        for (uint i = 0; i < poll.candidates.length; i++) {
            if (poll.votes[poll.candidates[i]] > maxVotes) {
                maxVotes = poll.votes[poll.candidates[i]];
                winningCandidate = poll.candidates[i];
            }
        }

        poll.winner = winningCandidate;
        poll.ended = true;

        // if (winningCandidate != address(0)) {
        //     payable(winningCandidate).transfer(1);
        // }

        emit PollEnded(pollId, winningCandidate);
    }

    function getCandidates(uint pollId) public view returns (address[] memory) {
        require(polls[pollId].creator != address(0), "Poll does not exist");
        Poll storage poll = polls[pollId];
        return poll.candidates;
    }

    function getOngoingPolls() public view returns (SimplePoll[] memory) {
        uint count = 0;
        for (uint i = 1; i <= pollCount; i++) {
            if (!polls[i].ended) {
                count++;
            }
        }

        SimplePoll[] memory ongoingPolls = new SimplePoll[](count);
        uint index = 0;
        for (uint i = 1; i <= pollCount; i++) {
            if (!polls[i].ended) {
                ongoingPolls[index] = SimplePoll({
                    pollId: i,
                    description: polls[i].description,
                    creator: polls[i].creator
                });
                index++;
            }
        }

        return ongoingPolls;
    }
}
