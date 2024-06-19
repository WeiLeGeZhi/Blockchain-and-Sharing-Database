// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

contract Ballot {
    struct Voter {
        bool joined;  // if true, that person already joined as a candidate
        bool voted;   // if true, that person already voted
        uint vote;    // index of the voted candidate
    }

    struct Poll {
        string description;
        address creator;
        uint endTime;
        uint reward;
        bool ended;
        address winner;
        address[] candidates;
        mapping(address => Voter) voters;
        mapping(address => uint) votes;
        uint voterCount; // Track number of voters
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

    event PollCreated(uint pollId, string description, uint endTime, uint reward);
    event CandidateJoined(uint pollId, address candidate);
    event Voted(uint pollId, address voter, address candidate);
    event PollEnded(uint pollId, address winner);

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier pollExists(uint pollId) {
        require(polls[pollId].creator != address(0), "Poll does not exist");
        _;
    }

    modifier onlyPollCreator(uint pollId) {
        require(polls[pollId].creator == msg.sender, "Only poll creator can call this function");
        _;
    }

    function createPoll(string memory description) public payable {
        pollCount++;
        uint duration = 10000;
        Poll storage poll = polls[pollCount];
        poll.description = description;
        poll.creator = msg.sender;
        poll.endTime = block.timestamp + duration;
        poll.reward = 1;
        poll.ended = false;
        poll.voterCount = 0;
        poll.winner = address(0);

        emit PollCreated(pollCount, poll.description, poll.endTime, poll.reward);
    }

    function joinPoll(uint pollId) public payable pollExists(pollId) {
        Poll storage poll = polls[pollId];
        require(!poll.ended, "Poll has ended");
        require(!poll.voters[msg.sender].joined, "Already joined");

        poll.voters[msg.sender].joined = true;
        poll.candidates.push(msg.sender);

        // payable(poll.creator).transfer(1);

        emit CandidateJoined(pollId, msg.sender);
    }

    function vote(uint pollId, address candidate) public pollExists(pollId) {
        Poll storage poll = polls[pollId];
        require(!poll.ended, "Poll has ended");
        require(!poll.voters[msg.sender].voted, "Already voted");
        require(!poll.voters[msg.sender].joined, "Candidates cannot vote for themselves");
        require(poll.voters[candidate].joined, "Candidate not joined");

        poll.voters[msg.sender].voted = true;
        poll.votes[candidate]++;
        poll.voterCount++;

        emit Voted(pollId, msg.sender, candidate);

        if (poll.voterCount >= 5) {
            endPoll(pollId);
        }
    }

    function endPoll(uint pollId) internal pollExists(pollId) {
        Poll storage poll = polls[pollId];
        require(!poll.ended, "Poll has already ended");

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
        //     payable(winningCandidate).transfer(poll.reward);
        // }

        emit PollEnded(pollId, winningCandidate);
    }

    function getCandidates(uint pollId) public view pollExists(pollId) returns (address[] memory) {
        Poll storage poll = polls[pollId];
        return poll.candidates;
    }

    function getPolls() public view returns (string memory) {
        bytes memory pollData;

        for (uint i = 1; i <= pollCount; i++) {
            Poll storage poll = polls[i];
            bytes memory candidatesData;
            for (uint j = 0; j < poll.candidates.length; j++) {
                candidatesData = abi.encodePacked(candidatesData, poll.candidates[j], ",");
            }
            candidatesData = abi.encodePacked("[", candidatesData, "]");

            pollData = abi.encodePacked(
                pollData,
                '{"description":"', poll.description,
                '","ended":', poll.ended ? "true" : "false",
                ',"winner":"', poll.winner,
                '","candidates":', candidatesData,
                '},'
            );
        }

        return string(abi.encodePacked("[", pollData, "]"));
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

    function getPollNum() public view returns (uint) {
        return pollCount;
    }

    function getWinner(uint pollId) public view pollExists(pollId) returns (address winner) {
        Poll storage poll = polls[pollId];
        return poll.winner;
    }
}
